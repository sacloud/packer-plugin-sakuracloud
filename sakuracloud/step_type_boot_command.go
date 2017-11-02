package sakuracloud

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/mitchellh/go-vnc"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/template/interpolate"
	"github.com/sacloud/libsacloud/sacloud"
)

const KeyLeftShift uint32 = 0xFFE1

type bootCommandTemplateData struct {
	ServerIP       string
	DefaultRoute   string
	NetworkMaskLen string
	DNS1           string
	DNS2           string
	PublicKey      string
	PrivateKey     string
}

// This step "types" the boot command into the VM over VNC.
//
// Uses:
//   config *config
//   ui     packer.Ui
//   vnc    *sacloud.VNCProxyResponse
//
// Produces:
//   <nothing>
type stepTypeBootCommand struct {
	Ctx   interpolate.Context
	Debug bool
}

func (s *stepTypeBootCommand) Run(state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(Config)
	if len(config.BootCommand) == 0 {
		return multistep.ActionContinue
	}

	ui := state.Get("ui").(packer.Ui)
	stepStartMsg(ui, s.Debug, "Type BootCommand")
	vncCredentials := state.Get("vnc").(*sacloud.VNCProxyResponse)

	serverIP := state.Get("server_ip").(string)
	defaultRoute := state.Get("default_route").(string)
	maskLen := state.Get("network_mask_len").(int)
	dns1 := state.Get("dns1").(string)
	dns2 := state.Get("dns2").(string)
	publicKey := state.Get("ssh_public_key").(string)
	privateKey := state.Get("ssh_private_key").(string)

	// Connect to VNC
	ui.Say("\tConnecting to VM via VNC")
	host := vncCredentials.ActualHost()
	nc, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, vncCredentials.Port))
	if err != nil {
		err := fmt.Errorf("Error connecting to VNC: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	defer nc.Close()

	auth := []vnc.ClientAuth{&vnc.PasswordAuth{Password: vncCredentials.Password}}
	c, err := vnc.Client(nc, &vnc.ClientConfig{Auth: auth, Exclusive: false})
	if err != nil {
		err := fmt.Errorf("Error handshaking with VNC: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	defer c.Close()

	log.Printf("Connected to VNC desktop: %s", c.DesktopName)

	s.Ctx.Data = &bootCommandTemplateData{
		ServerIP:       serverIP,
		DefaultRoute:   defaultRoute,
		NetworkMaskLen: fmt.Sprintf("%d", maskLen),
		DNS1:           dns1,
		DNS2:           dns2,
		PublicKey:      publicKey,
		PrivateKey:     privateKey,
	}

	ui.Say("\tTyping the boot command over VNC...")
	for _, command := range config.BootCommand {
		command, err := interpolate.Render(command, &s.Ctx)
		if err != nil {
			err := fmt.Errorf("Error preparing boot command: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}

		// Check for interrupts between typing things so we can cancel
		// since this isn't the fastest thing.
		if _, ok := state.GetOk(multistep.StateCancelled); ok {
			return multistep.ActionHalt
		}

		vncSendString(c, command, config.UseUSKeyboard)
	}

	stepEndMsg(ui, s.Debug, "BootWait")
	return multistep.ActionContinue
}

func (*stepTypeBootCommand) Cleanup(multistep.StateBag) {}

func vncSendString(c *vnc.ClientConn, original string, useUSKeyboard bool) {
	// Scancodes reference: https://github.com/qemu/qemu/blob/master/ui/vnc_keysym.h
	special := make(map[string]uint32)
	special["<bs>"] = 0xFF08
	special["<del>"] = 0xFFFF
	special["<enter>"] = 0xFF0D
	special["<esc>"] = 0xFF1B
	special["<f1>"] = 0xFFBE
	special["<f2>"] = 0xFFBF
	special["<f3>"] = 0xFFC0
	special["<f4>"] = 0xFFC1
	special["<f5>"] = 0xFFC2
	special["<f6>"] = 0xFFC3
	special["<f7>"] = 0xFFC4
	special["<f8>"] = 0xFFC5
	special["<f9>"] = 0xFFC6
	special["<f10>"] = 0xFFC7
	special["<f11>"] = 0xFFC8
	special["<f12>"] = 0xFFC9
	special["<return>"] = 0xFF0D
	special["<tab>"] = 0xFF09
	special["<up>"] = 0xFF52
	special["<down>"] = 0xFF54
	special["<left>"] = 0xFF51
	special["<right>"] = 0xFF53
	special["<spacebar>"] = 0x020
	special["<insert>"] = 0xFF63
	special["<home>"] = 0xFF50
	special["<end>"] = 0xFF57
	special["<pageUp>"] = 0xFF55
	special["<pageDown>"] = 0xFF56
	special["<leftAlt>"] = 0xFFE9
	special["<leftCtrl>"] = 0xFFE3
	special["<leftShift>"] = 0xFFE1
	special["<rightAlt>"] = 0xFFEA
	special["<rightCtrl>"] = 0xFFE4
	special["<rightShift>"] = 0xFFE2
	special["<leftWin>"] = 0xFF5B
	special["<rightWin>"] = 0xFF5C

	shiftedChars := "!\"#$%&'()=~|{`+*}<>?"
	if useUSKeyboard {
		shiftedChars = "~!@#$%^&*()_+{}|:\"<>?"
	}

	// TODO(mitchellh): Ripe for optimizations of some point, perhaps.
	for len(original) > 0 {
		var keyCode uint32
		keyShift := false

		if strings.HasPrefix(original, "<leftAltOn>") {
			keyCode = special["<leftAlt>"]
			original = original[len("<leftAltOn>"):]
			log.Printf("Special code '<leftAltOn>' found, replacing with: %d", keyCode)

			c.KeyEvent(keyCode, true)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<leftCtrlOn>") {
			keyCode = special["<leftCtrl>"]
			original = original[len("<leftCtrlOn>"):]
			log.Printf("Special code '<leftCtrlOn>' found, replacing with: %d", keyCode)

			c.KeyEvent(keyCode, true)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<leftShiftOn>") {
			keyCode = special["<leftShift>"]
			original = original[len("<leftShiftOn>"):]
			log.Printf("Special code '<leftShiftOn>' found, replacing with: %d", keyCode)

			c.KeyEvent(keyCode, true)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<leftAltOff>") {
			keyCode = special["<leftAlt>"]
			original = original[len("<leftAltOff>"):]
			log.Printf("Special code '<leftAltOff>' found, replacing with: %d", keyCode)

			c.KeyEvent(keyCode, false)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<leftCtrlOff>") {
			keyCode = special["<leftCtrl>"]
			original = original[len("<leftCtrlOff>"):]
			log.Printf("Special code '<leftCtrlOff>' found, replacing with: %d", keyCode)

			c.KeyEvent(keyCode, false)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<leftShiftOff>") {
			keyCode = special["<leftShift>"]
			original = original[len("<leftShiftOff>"):]
			log.Printf("Special code '<leftShiftOff>' found, replacing with: %d", keyCode)

			c.KeyEvent(keyCode, false)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<rightAltOn>") {
			keyCode = special["<rightAlt>"]
			original = original[len("<rightAltOn>"):]
			log.Printf("Special code '<rightAltOn>' found, replacing with: %d", keyCode)

			c.KeyEvent(keyCode, true)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<rightCtrlOn>") {
			keyCode = special["<rightCtrl>"]
			original = original[len("<rightCtrlOn>"):]
			log.Printf("Special code '<rightCtrlOn>' found, replacing with: %d", keyCode)

			c.KeyEvent(keyCode, true)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<rightShiftOn>") {
			keyCode = special["<rightShift>"]
			original = original[len("<rightShiftOn>"):]
			log.Printf("Special code '<rightShiftOn>' found, replacing with: %d", keyCode)

			c.KeyEvent(keyCode, true)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<rightAltOff>") {
			keyCode = special["<rightAlt>"]
			original = original[len("<rightAltOff>"):]
			log.Printf("Special code '<rightAltOff>' found, replacing with: %d", keyCode)

			c.KeyEvent(keyCode, false)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<rightCtrlOff>") {
			keyCode = special["<rightCtrl>"]
			original = original[len("<rightCtrlOff>"):]
			log.Printf("Special code '<rightCtrlOff>' found, replacing with: %d", keyCode)

			c.KeyEvent(keyCode, false)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<rightShiftOff>") {
			keyCode = special["<rightShift>"]
			original = original[len("<rightShiftOff>"):]
			log.Printf("Special code '<rightShiftOff>' found, replacing with: %d", keyCode)

			c.KeyEvent(keyCode, false)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<leftWinOn>") {
			keyCode = special["<leftWin>"]
			original = original[len("<leftWinOn>"):]
			log.Printf("Special code '<leftWinOn>' found, replacing with: %d", keyCode)

			c.KeyEvent(keyCode, true)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<leftWinOff>") {
			keyCode = special["<leftWin>"]
			original = original[len("<leftWinOff>"):]
			log.Printf("Special code '<leftWinOff>' found, replacing with: %d", keyCode)

			c.KeyEvent(keyCode, false)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}
		if strings.HasPrefix(original, "<rightWinOn>") {
			keyCode = special["<rightWin>"]
			original = original[len("<rightWinOn>"):]
			log.Printf("Special code '<rightWinOn>' found, replacing with: %d", keyCode)

			c.KeyEvent(keyCode, true)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<rightWinOff>") {
			keyCode = special["<rightWin>"]
			original = original[len("<rightWinOff>"):]
			log.Printf("Special code '<rightWinOff>' found, replacing with: %d", keyCode)

			c.KeyEvent(keyCode, false)
			time.Sleep(time.Second / 10)

			// qemu is picky, so no matter what, wait a small period
			time.Sleep(100 * time.Millisecond)

			continue
		}

		if strings.HasPrefix(original, "<wait>") {
			log.Printf("Special code '<wait>' found, sleeping one second")
			time.Sleep(1 * time.Second)
			original = original[len("<wait>"):]
			continue
		}

		if strings.HasPrefix(original, "<wait5>") {
			log.Printf("Special code '<wait5>' found, sleeping 5 seconds")
			time.Sleep(5 * time.Second)
			original = original[len("<wait5>"):]
			continue
		}

		if strings.HasPrefix(original, "<wait10>") {
			log.Printf("Special code '<wait10>' found, sleeping 10 seconds")
			time.Sleep(10 * time.Second)
			original = original[len("<wait10>"):]
			continue
		}

		if strings.HasPrefix(original, "<wait") && strings.HasSuffix(original, ">") {
			re := regexp.MustCompile(`<wait([0-9hms]+)>$`)
			dstr := re.FindStringSubmatch(original)
			if len(dstr) > 1 {
				log.Printf("Special code %s found, sleeping", dstr[0])
				if dt, err := time.ParseDuration(dstr[1]); err == nil {
					time.Sleep(dt)
					original = original[len(dstr[0]):]
					continue
				}
			}
		}

		for specialCode, specialValue := range special {
			if strings.HasPrefix(original, specialCode) {
				log.Printf("Special code '%s' found, replacing with: %d", specialCode, specialValue)
				keyCode = specialValue
				original = original[len(specialCode):]
				break
			}
		}

		if keyCode == 0 {
			r, size := utf8.DecodeRuneInString(original)
			original = original[size:]
			keyCode = uint32(r)
			keyShift = unicode.IsUpper(r) || strings.ContainsRune(shiftedChars, r)

			log.Printf("Sending char '%c', code %d, shift %v", r, keyCode, keyShift)
		}

		if keyShift {
			c.KeyEvent(KeyLeftShift, true)
		}

		c.KeyEvent(keyCode, true)
		time.Sleep(time.Second / 10)
		c.KeyEvent(keyCode, false)
		time.Sleep(time.Second / 10)

		if keyShift {
			c.KeyEvent(KeyLeftShift, false)
		}

		// qemu is picky, so no matter what, wait a small period
		time.Sleep(100 * time.Millisecond)
	}
}
