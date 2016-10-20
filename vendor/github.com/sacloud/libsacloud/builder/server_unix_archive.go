package builder

import "github.com/sacloud/libsacloud/sacloud"

// PublicArchiveUnixServerBuilder Linux(Unix)系パブリックアーカイブを利用して構築を行うサーバービルダー
//
// 基本的なディスク設定やディスクの追加、ディスクの修正機能に対応しています。
type PublicArchiveUnixServerBuilder struct {
	*serverBuilder
}

/*---------------------------------------------------------
  common properties
---------------------------------------------------------*/

// GetServerName サーバー名 取得
func (b *PublicArchiveUnixServerBuilder) GetServerName() string {
	return b.serverName
}

// SetServerName サーバー名 設定
func (b *PublicArchiveUnixServerBuilder) SetServerName(serverName string) *PublicArchiveUnixServerBuilder {
	b.serverName = serverName
	return b
}

// GetCore CPUコア数 取得
func (b *PublicArchiveUnixServerBuilder) GetCore() int {
	return b.core
}

// SetCore CPUコア数 設定
func (b *PublicArchiveUnixServerBuilder) SetCore(core int) *PublicArchiveUnixServerBuilder {
	b.core = core
	return b
}

// GetMemory メモリサイズ(GB単位) 取得
func (b *PublicArchiveUnixServerBuilder) GetMemory() int {
	return b.memory
}

// SetMemory メモリサイズ(GB単位) 設定
func (b *PublicArchiveUnixServerBuilder) SetMemory(memory int) *PublicArchiveUnixServerBuilder {
	b.memory = memory
	return b
}

// IsUseVirtIONetPCI NIC準仮装化モード(virtio)利用フラグ 取得
func (b *PublicArchiveUnixServerBuilder) IsUseVirtIONetPCI() bool {
	return b.useVirtIONetPCI
}

// SetUseVirtIONetPCI NIC準仮装化モード(virtio)利用フラグ 設定
func (b *PublicArchiveUnixServerBuilder) SetUseVirtIONetPCI(useVirtIONetPCI bool) *PublicArchiveUnixServerBuilder {
	b.useVirtIONetPCI = useVirtIONetPCI
	return b
}

// GetDescription 説明 取得
func (b *PublicArchiveUnixServerBuilder) GetDescription() string {
	return b.description
}

// SetDescription 説明 設定
func (b *PublicArchiveUnixServerBuilder) SetDescription(description string) *PublicArchiveUnixServerBuilder {
	b.description = description
	return b
}

// GetIconID アイコンID 取得
func (b *PublicArchiveUnixServerBuilder) GetIconID() int64 {
	return b.iconID
}

// SetIconID アイコンID 設定
func (b *PublicArchiveUnixServerBuilder) SetIconID(iconID int64) *PublicArchiveUnixServerBuilder {
	b.iconID = iconID
	return b
}

// IsBootAfterCreate サーバー作成後すぐに起動フラグ 取得
func (b *PublicArchiveUnixServerBuilder) IsBootAfterCreate() bool {
	return b.bootAfterCreate
}

// SetBootAfterCreate サーバー作成後すぐに起動フラグ 設定
func (b *PublicArchiveUnixServerBuilder) SetBootAfterCreate(bootAfterCreate bool) *PublicArchiveUnixServerBuilder {
	b.bootAfterCreate = bootAfterCreate
	return b
}

// GetTags タグ 取得
func (b *PublicArchiveUnixServerBuilder) GetTags() []string {
	return b.Tags
}

// SetTags タグ 設定
func (b *PublicArchiveUnixServerBuilder) SetTags(tags []string) *PublicArchiveUnixServerBuilder {
	b.Tags = tags
	return b
}

/*---------------------------------------------------------
  for nic functioms
---------------------------------------------------------*/

// ClearNICConnections NIC接続設定 クリア
func (b *PublicArchiveUnixServerBuilder) ClearNICConnections() *PublicArchiveUnixServerBuilder {
	b.nicConnections = nil
	b.disk.SetIPAddress("")
	b.disk.SetNetworkMaskLen(0)
	b.disk.SetDefaultRoute("")
	return b
}

// AddPublicNWConnectedNIC 共有セグメントへの接続追加(注:共有セグメントはeth0のみ接続可能)
func (b *PublicArchiveUnixServerBuilder) AddPublicNWConnectedNIC() *PublicArchiveUnixServerBuilder {
	b.nicConnections = append(b.nicConnections, "shared")
	return b
}

// AddExistsSwitchConnectedNIC スイッチ or ルーター+スイッチへの接続追加(注:ルーター+スイッチはeth0のみ接続可能)
func (b *PublicArchiveUnixServerBuilder) AddExistsSwitchConnectedNIC(switchID string, ipaddress string, networkMaskLen int, defaultRoute string) *PublicArchiveUnixServerBuilder {
	b.nicConnections = append(b.nicConnections, switchID)
	b.disk.SetIPAddress(ipaddress)
	b.disk.SetNetworkMaskLen(networkMaskLen)
	b.disk.SetDefaultRoute(defaultRoute)
	return b
}

// AddDisconnectedNIC 切断されたNIC追加
func (b *PublicArchiveUnixServerBuilder) AddDisconnectedNIC() *PublicArchiveUnixServerBuilder {
	b.nicConnections = append(b.nicConnections, "")
	return b
}

// GetISOImageID ISOイメージ(CDROM)ID 取得
func (b *PublicArchiveUnixServerBuilder) GetISOImageID() int64 {
	return b.isoImageID
}

// SetISOImageID ISOイメージ(CDROM)ID 設定
func (b *PublicArchiveUnixServerBuilder) SetISOImageID(id int64) *PublicArchiveUnixServerBuilder {
	b.isoImageID = id
	return b
}

/*---------------------------------------------------------
  for disk properties
---------------------------------------------------------*/

// GetDiskSize ディスクサイズ(GB単位) 取得
func (b *PublicArchiveUnixServerBuilder) GetDiskSize() int {
	return b.disk.GetSize()
}

// SetDiskSize ディスクサイズ(GB単位) 設定
func (b *PublicArchiveUnixServerBuilder) SetDiskSize(diskSize int) *PublicArchiveUnixServerBuilder {
	b.disk.SetSize(diskSize)
	return b
}

// GetDistantFrom ストレージ隔離対象ディスク 取得
func (b *PublicArchiveUnixServerBuilder) GetDistantFrom() []int64 {
	return b.disk.GetDistantFrom()
}

// SetDistantFrom ストレージ隔離対象ディスク 設定
func (b *PublicArchiveUnixServerBuilder) SetDistantFrom(distantFrom []int64) *PublicArchiveUnixServerBuilder {
	b.disk.SetDistantFrom(distantFrom)
	return b
}

// AddDistantFrom ストレージ隔離対象ディスク 追加
func (b *PublicArchiveUnixServerBuilder) AddDistantFrom(diskID int64) *PublicArchiveUnixServerBuilder {
	b.disk.AddDistantFrom(diskID)
	return b
}

// ClearDistantFrom ストレージ隔離対象ディスク クリア
func (b *PublicArchiveUnixServerBuilder) ClearDistantFrom() *PublicArchiveUnixServerBuilder {
	b.disk.ClearDistantFrom()
	return b
}

// GetDiskPlanID ディスクプラン(SSD/HDD) 取得
func (b *PublicArchiveUnixServerBuilder) GetDiskPlanID() sacloud.DiskPlanID {
	return b.disk.GetPlanID()
}

// SetDiskPlanID ディスクプラン(SSD/HDD) 設定
func (b *PublicArchiveUnixServerBuilder) SetDiskPlanID(diskPlanID sacloud.DiskPlanID) *PublicArchiveUnixServerBuilder {
	b.disk.SetPlanID(diskPlanID)
	return b
}

// GetDiskConnection ディスク接続方法(VirtIO/IDE) 取得
func (b *PublicArchiveUnixServerBuilder) GetDiskConnection() sacloud.EDiskConnection {
	return b.disk.GetConnection()
}

// SetDiskConnection ディスク接続方法(VirtIO/IDE) 設定
func (b *PublicArchiveUnixServerBuilder) SetDiskConnection(diskConnection sacloud.EDiskConnection) *PublicArchiveUnixServerBuilder {
	b.disk.SetConnection(diskConnection)
	return b
}

/*---------------------------------------------------------
  for disk edit properties
---------------------------------------------------------*/

// GetSourceArchiveID ソースアーカイブID 取得
func (b *PublicArchiveUnixServerBuilder) GetSourceArchiveID() int64 {
	return b.disk.GetSourceArchiveID()
}

// GetSourceDiskID ソースディスクID 設定
func (b *PublicArchiveUnixServerBuilder) GetSourceDiskID() int64 {
	return b.disk.GetSourceDiskID()
}

// GetPassword パスワード 取得
func (b *PublicArchiveUnixServerBuilder) GetPassword() string {
	return b.disk.GetPassword()
}

// SetPassword パスワード 設定
func (b *PublicArchiveUnixServerBuilder) SetPassword(password string) *PublicArchiveUnixServerBuilder {
	b.disk.SetPassword(password)
	return b
}

// GetHostName ホスト名 取得
func (b *PublicArchiveUnixServerBuilder) GetHostName() string {
	return b.disk.GetHostName()
}

// SetHostName ホスト名 設定
func (b *PublicArchiveUnixServerBuilder) SetHostName(hostName string) *PublicArchiveUnixServerBuilder {
	b.disk.SetHostName(hostName)
	return b
}

// IsDisablePWAuth パスワード認証無効化フラグ 取得
func (b *PublicArchiveUnixServerBuilder) IsDisablePWAuth() bool {
	return b.disk.IsDisablePWAuth()
}

// SetDisablePWAuth パスワード認証無効化フラグ 設定
func (b *PublicArchiveUnixServerBuilder) SetDisablePWAuth(disable bool) *PublicArchiveUnixServerBuilder {
	b.disk.SetDisablePWAuth(disable)
	return b
}

// AddSSHKey 公開鍵 追加
func (b *PublicArchiveUnixServerBuilder) AddSSHKey(sshKey string) *PublicArchiveUnixServerBuilder {
	b.disk.AddSSHKey(sshKey)
	return b
}

// ClearSSHKey 公開鍵 クリア
func (b *PublicArchiveUnixServerBuilder) ClearSSHKey() *PublicArchiveUnixServerBuilder {
	b.disk.ClearSSHKey()
	return b
}

// GetSSHKeys 公開鍵 取得
func (b *PublicArchiveUnixServerBuilder) GetSSHKeys() []string {
	return b.disk.GetSSHKeys()
}

// AddSSHKeyID 公開鍵ID 追加
func (b *PublicArchiveUnixServerBuilder) AddSSHKeyID(sshKeyID int64) *PublicArchiveUnixServerBuilder {
	b.disk.AddSSHKeyID(sshKeyID)
	return b
}

// ClearSSHKeyIDs 公開鍵ID クリア
func (b *PublicArchiveUnixServerBuilder) ClearSSHKeyIDs() *PublicArchiveUnixServerBuilder {
	b.disk.ClearSSHKeyIDs()
	return b
}

// GetSSHKeyIds 公開鍵ID 取得
func (b *PublicArchiveUnixServerBuilder) GetSSHKeyIds() []int64 {
	return b.disk.GetSSHKeyIds()
}

// AddNote スタートアップスクリプト 追加
func (b *PublicArchiveUnixServerBuilder) AddNote(note string) *PublicArchiveUnixServerBuilder {
	b.disk.AddNote(note)
	return b
}

// ClearNotes スタートアップスクリプト クリア
func (b *PublicArchiveUnixServerBuilder) ClearNotes() *PublicArchiveUnixServerBuilder {
	b.disk.ClearNotes()
	return b
}

// GetNotes スタートアップスクリプト 取得
func (b *PublicArchiveUnixServerBuilder) GetNotes() []string {
	return b.disk.GetNotes()
}

// AddNoteID スタートアップスクリプト 追加
func (b *PublicArchiveUnixServerBuilder) AddNoteID(noteID int64) *PublicArchiveUnixServerBuilder {
	b.disk.AddNoteID(noteID)
	return b
}

// ClearNoteIDs スタートアップスクリプト クリア
func (b *PublicArchiveUnixServerBuilder) ClearNoteIDs() *PublicArchiveUnixServerBuilder {
	b.disk.ClearNoteIDs()
	return b
}

// GetNoteIDs スタートアップスクリプトID 取得
func (b *PublicArchiveUnixServerBuilder) GetNoteIDs() []int64 {
	return b.disk.GetNoteIDs()
}

// IsSSHKeysEphemeral ディスク作成後の公開鍵削除フラグ 取得
func (b *PublicArchiveUnixServerBuilder) IsSSHKeysEphemeral() bool {
	return b.disk.IsSSHKeysEphemeral()
}

// SetSSHKeysEphemeral ディスク作成後の公開鍵削除フラグ 設定
func (b *PublicArchiveUnixServerBuilder) SetSSHKeysEphemeral(isEphemeral bool) *PublicArchiveUnixServerBuilder {
	b.disk.SetSSHKeysEphemeral(isEphemeral)
	return b
}

// IsNotesEphemeral ディスク作成後のスタートアップスクリプト削除フラグ 取得
func (b *PublicArchiveUnixServerBuilder) IsNotesEphemeral() bool {
	return b.disk.IsNotesEphemeral()
}

// SetNotesEphemeral ディスク作成後のスタートアップスクリプト削除フラグ 設定
func (b *PublicArchiveUnixServerBuilder) SetNotesEphemeral(isEphemeral bool) *PublicArchiveUnixServerBuilder {
	b.disk.SetNotesEphemeral(isEphemeral)
	return b
}

/*---------------------------------------------------------
  for event handler
---------------------------------------------------------*/

// SetEventHandler イベントハンドラ 設定
func (b *PublicArchiveUnixServerBuilder) SetEventHandler(event ServerBuildEvents, handler ServerBuildEventHandler) *PublicArchiveUnixServerBuilder {
	b.buildEventHandlers[event] = handler
	return b
}

// ClearEventHandler イベントハンドラ クリア
func (b *PublicArchiveUnixServerBuilder) ClearEventHandler(event ServerBuildEvents) *PublicArchiveUnixServerBuilder {
	delete(b.buildEventHandlers, event)
	return b
}

// GetEventHandler イベントハンドラ 取得
func (b *PublicArchiveUnixServerBuilder) GetEventHandler(event ServerBuildEvents) *ServerBuildEventHandler {
	if handler, ok := b.buildEventHandlers[event]; ok {
		return &handler
	}
	return nil
}

// SetDiskEventHandler ディスクイベントハンドラ 設定
func (b *PublicArchiveUnixServerBuilder) SetDiskEventHandler(event DiskBuildEvents, handler DiskBuildEventHandler) *PublicArchiveUnixServerBuilder {
	b.disk.SetEventHandler(event, handler)
	return b
}

// ClearDiskEventHandler ディスクイベントハンドラ クリア
func (b *PublicArchiveUnixServerBuilder) ClearDiskEventHandler(event DiskBuildEvents) *PublicArchiveUnixServerBuilder {
	b.disk.ClearEventHandler(event)
	return b
}

// GetDiskEventHandler ディスクイベントハンドラ 取得
func (b *PublicArchiveUnixServerBuilder) GetDiskEventHandler(event DiskBuildEvents) *DiskBuildEventHandler {
	return b.disk.GetEventHandler(event)
}

/*---------------------------------------------------------
  for additional disks
---------------------------------------------------------*/

// AddAdditionalDisk 追加ディスク 追加
func (b *PublicArchiveUnixServerBuilder) AddAdditionalDisk(diskBuilder *DiskBuilder) *PublicArchiveUnixServerBuilder {
	b.additionalDisks = append(b.additionalDisks, diskBuilder)
	return b
}

// ClearAdditionalDisks 追加ディスク クリア
func (b *PublicArchiveUnixServerBuilder) ClearAdditionalDisks() *PublicArchiveUnixServerBuilder {
	b.additionalDisks = []*DiskBuilder{}
	return b
}

// GetAdditionalDisks 追加ディスク 取得
func (b *PublicArchiveUnixServerBuilder) GetAdditionalDisks() []*DiskBuilder {
	return b.additionalDisks
}
