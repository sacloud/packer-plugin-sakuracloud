package builder

import "github.com/sacloud/libsacloud/sacloud"

// CommonServerBuilder 既存のアーカイブ or ディスクを利用して構築を行うサーバービルダー
//
// 基本的なディスク設定やディスクの追加、ディスクの修正機能に対応しています。
// ただし、ディスクの修正機能はソースアーカイブ or ソースディスクが対応していない場合は
// さくらのクラウドAPIコール時にエラーとなるため、適切にハンドリングするように実装する必要があります。
type CommonServerBuilder struct {
	*serverBuilder
}

/*---------------------------------------------------------
  common properties
---------------------------------------------------------*/

// GetServerName サーバー名 取得
func (b *CommonServerBuilder) GetServerName() string {
	return b.serverName
}

// SetServerName サーバー名 設定
func (b *CommonServerBuilder) SetServerName(serverName string) *CommonServerBuilder {
	b.serverName = serverName
	return b
}

// GetCore CPUコア数 取得
func (b *CommonServerBuilder) GetCore() int {
	return b.core
}

// SetCore CPUコア数 設定
func (b *CommonServerBuilder) SetCore(core int) *CommonServerBuilder {
	b.core = core
	return b
}

// GetMemory メモリサイズ(GB単位) 取得
func (b *CommonServerBuilder) GetMemory() int {
	return b.memory
}

// SetMemory メモリサイズ(GB単位) 設定
func (b *CommonServerBuilder) SetMemory(memory int) *CommonServerBuilder {
	b.memory = memory
	return b
}

// IsUseVirtIONetPCI NIC準仮装化モード(virtio)利用フラグ 取得
func (b *CommonServerBuilder) IsUseVirtIONetPCI() bool {
	return b.useVirtIONetPCI
}

// SetUseVirtIONetPCI NIC準仮装化モード(virtio)利用フラグ 設定
func (b *CommonServerBuilder) SetUseVirtIONetPCI(useVirtIONetPCI bool) *CommonServerBuilder {
	b.useVirtIONetPCI = useVirtIONetPCI
	return b
}

// GetDescription 説明 取得
func (b *CommonServerBuilder) GetDescription() string {
	return b.description
}

// SetDescription 説明 設定
func (b *CommonServerBuilder) SetDescription(description string) *CommonServerBuilder {
	b.description = description
	return b
}

// GetIconID アイコンID 取得
func (b *CommonServerBuilder) GetIconID() int64 {
	return b.iconID
}

// SetIconID アイコンID 設定
func (b *CommonServerBuilder) SetIconID(iconID int64) *CommonServerBuilder {
	b.iconID = iconID
	return b
}

// IsBootAfterCreate サーバー作成後すぐに起動フラグ 取得
func (b *CommonServerBuilder) IsBootAfterCreate() bool {
	return b.bootAfterCreate
}

// SetBootAfterCreate サーバー作成後すぐに起動フラグ 設定
func (b *CommonServerBuilder) SetBootAfterCreate(bootAfterCreate bool) *CommonServerBuilder {
	b.bootAfterCreate = bootAfterCreate
	return b
}

// GetTags タグ 取得
func (b *CommonServerBuilder) GetTags() []string {
	return b.Tags
}

// SetTags タグ 設定
func (b *CommonServerBuilder) SetTags(tags []string) *CommonServerBuilder {
	b.Tags = tags
	return b
}

/*---------------------------------------------------------
  for nic functioms
---------------------------------------------------------*/

// ClearNICConnections NIC接続設定 クリア
func (b *CommonServerBuilder) ClearNICConnections() *CommonServerBuilder {
	b.nicConnections = nil
	b.disk.SetIPAddress("")
	b.disk.SetNetworkMaskLen(0)
	b.disk.SetDefaultRoute("")
	return b
}

// AddPublicNWConnectedNIC 共有セグメントへの接続追加(注:共有セグメントはeth0のみ接続可能)
func (b *CommonServerBuilder) AddPublicNWConnectedNIC() *CommonServerBuilder {
	b.nicConnections = append(b.nicConnections, "shared")
	return b
}

// AddExistsSwitchConnectedNIC スイッチ or ルーター+スイッチへの接続追加(注:ルーター+スイッチはeth0のみ接続可能)
func (b *CommonServerBuilder) AddExistsSwitchConnectedNIC(switchID string, ipaddress string, networkMaskLen int, defaultRoute string) *CommonServerBuilder {
	b.nicConnections = append(b.nicConnections, switchID)
	b.disk.SetIPAddress(ipaddress)
	b.disk.SetNetworkMaskLen(networkMaskLen)
	b.disk.SetDefaultRoute(defaultRoute)
	return b
}

// AddDisconnectedNIC 切断されたNIC追加
func (b *CommonServerBuilder) AddDisconnectedNIC() *CommonServerBuilder {
	b.nicConnections = append(b.nicConnections, "")
	return b
}

// GetISOImageID ISOイメージ(CDROM)ID 取得
func (b *CommonServerBuilder) GetISOImageID() int64 {
	return b.isoImageID
}

// SetISOImageID ISOイメージ(CDROM)ID 設定
func (b *CommonServerBuilder) SetISOImageID(id int64) *CommonServerBuilder {
	b.isoImageID = id
	return b
}

/*---------------------------------------------------------
  for disk properties
---------------------------------------------------------*/

// GetDiskSize ディスクサイズ(GB単位) 取得
func (b *CommonServerBuilder) GetDiskSize() int {
	return b.disk.GetSize()
}

// SetDiskSize ディスクサイズ(GB単位) 設定
func (b *CommonServerBuilder) SetDiskSize(diskSize int) *CommonServerBuilder {
	b.disk.SetSize(diskSize)
	return b
}

// GetDistantFrom ストレージ隔離対象ディスク 取得
func (b *CommonServerBuilder) GetDistantFrom() []int64 {
	return b.disk.GetDistantFrom()
}

// SetDistantFrom ストレージ隔離対象ディスク 設定
func (b *CommonServerBuilder) SetDistantFrom(distantFrom []int64) *CommonServerBuilder {
	b.disk.SetDistantFrom(distantFrom)
	return b
}

// AddDistantFrom ストレージ隔離対象ディスク 追加
func (b *CommonServerBuilder) AddDistantFrom(diskID int64) *CommonServerBuilder {
	b.disk.AddDistantFrom(diskID)
	return b
}

// ClearDistantFrom ストレージ隔離対象ディスク クリア
func (b *CommonServerBuilder) ClearDistantFrom() *CommonServerBuilder {
	b.disk.ClearDistantFrom()
	return b
}

// GetDiskPlanID ディスクプラン(SSD/HDD) 取得
func (b *CommonServerBuilder) GetDiskPlanID() sacloud.DiskPlanID {
	return b.disk.GetPlanID()
}

// SetDiskPlanID ディスクプラン(SSD/HDD) 設定
func (b *CommonServerBuilder) SetDiskPlanID(diskPlanID sacloud.DiskPlanID) *CommonServerBuilder {
	b.disk.SetPlanID(diskPlanID)
	return b
}

// GetDiskConnection ディスク接続方法(VirtIO/IDE) 取得
func (b *CommonServerBuilder) GetDiskConnection() sacloud.EDiskConnection {
	return b.disk.GetConnection()
}

// SetDiskConnection ディスク接続方法(VirtIO/IDE) 設定
func (b *CommonServerBuilder) SetDiskConnection(diskConnection sacloud.EDiskConnection) *CommonServerBuilder {
	b.disk.SetConnection(diskConnection)
	return b
}

/*---------------------------------------------------------
  for disk edit properties
---------------------------------------------------------*/

// GetSourceArchiveID ソースアーカイブID 取得
func (b *CommonServerBuilder) GetSourceArchiveID() int64 {
	return b.disk.GetSourceArchiveID()
}

// GetSourceDiskID ソースディスクID 設定
func (b *CommonServerBuilder) GetSourceDiskID() int64 {
	return b.disk.GetSourceDiskID()
}

// GetPassword パスワード 取得
func (b *CommonServerBuilder) GetPassword() string {
	return b.disk.GetPassword()
}

// SetPassword パスワード 設定
func (b *CommonServerBuilder) SetPassword(password string) *CommonServerBuilder {
	b.disk.SetPassword(password)
	return b
}

// GetHostName ホスト名 取得
func (b *CommonServerBuilder) GetHostName() string {
	return b.disk.GetHostName()
}

// SetHostName ホスト名 設定
func (b *CommonServerBuilder) SetHostName(hostName string) *CommonServerBuilder {
	b.disk.SetHostName(hostName)
	return b
}

// IsDisablePWAuth パスワード認証無効化フラグ 取得
func (b *CommonServerBuilder) IsDisablePWAuth() bool {
	return b.disk.IsDisablePWAuth()
}

// SetDisablePWAuth パスワード認証無効化フラグ 設定
func (b *CommonServerBuilder) SetDisablePWAuth(disable bool) *CommonServerBuilder {
	b.disk.SetDisablePWAuth(disable)
	return b
}

// AddSSHKey 公開鍵 追加
func (b *CommonServerBuilder) AddSSHKey(sshKey string) *CommonServerBuilder {
	b.disk.AddSSHKey(sshKey)
	return b
}

// ClearSSHKey 公開鍵 クリア
func (b *CommonServerBuilder) ClearSSHKey() *CommonServerBuilder {
	b.disk.ClearSSHKey()
	return b
}

// GetSSHKeys 公開鍵 取得
func (b *CommonServerBuilder) GetSSHKeys() []string {
	return b.disk.GetSSHKeys()
}

// AddSSHKeyID 公開鍵ID 追加
func (b *CommonServerBuilder) AddSSHKeyID(sshKeyID int64) *CommonServerBuilder {
	b.disk.AddSSHKeyID(sshKeyID)
	return b
}

// ClearSSHKeyIDs 公開鍵ID クリア
func (b *CommonServerBuilder) ClearSSHKeyIDs() *CommonServerBuilder {
	b.disk.ClearSSHKeyIDs()
	return b
}

// GetSSHKeyIds 公開鍵ID 取得
func (b *CommonServerBuilder) GetSSHKeyIds() []int64 {
	return b.disk.GetSSHKeyIds()
}

// AddNote スタートアップスクリプト 追加
func (b *CommonServerBuilder) AddNote(note string) *CommonServerBuilder {
	b.disk.AddNote(note)
	return b
}

// ClearNotes スタートアップスクリプト クリア
func (b *CommonServerBuilder) ClearNotes() *CommonServerBuilder {
	b.disk.ClearNotes()
	return b
}

// GetNotes スタートアップスクリプト 取得
func (b *CommonServerBuilder) GetNotes() []string {
	return b.disk.GetNotes()
}

// AddNoteID スタートアップスクリプト 追加
func (b *CommonServerBuilder) AddNoteID(noteID int64) *CommonServerBuilder {
	b.disk.AddNoteID(noteID)
	return b
}

// ClearNoteIDs スタートアップスクリプト クリア
func (b *CommonServerBuilder) ClearNoteIDs() *CommonServerBuilder {
	b.disk.ClearNoteIDs()
	return b
}

// GetNoteIDs スタートアップスクリプトID 取得
func (b *CommonServerBuilder) GetNoteIDs() []int64 {
	return b.disk.GetNoteIDs()
}

// IsSSHKeysEphemeral ディスク作成後の公開鍵削除フラグ 取得
func (b *CommonServerBuilder) IsSSHKeysEphemeral() bool {
	return b.disk.IsSSHKeysEphemeral()
}

// SetSSHKeysEphemeral ディスク作成後の公開鍵削除フラグ 設定
func (b *CommonServerBuilder) SetSSHKeysEphemeral(isEphemeral bool) *CommonServerBuilder {
	b.disk.SetSSHKeysEphemeral(isEphemeral)
	return b
}

// IsNotesEphemeral ディスク作成後のスタートアップスクリプト削除フラグ 取得
func (b *CommonServerBuilder) IsNotesEphemeral() bool {
	return b.disk.IsNotesEphemeral()
}

// SetNotesEphemeral ディスク作成後のスタートアップスクリプト削除フラグ 設定
func (b *CommonServerBuilder) SetNotesEphemeral(isEphemeral bool) *CommonServerBuilder {
	b.disk.SetNotesEphemeral(isEphemeral)
	return b
}

/*---------------------------------------------------------
  for event handler
---------------------------------------------------------*/

// SetEventHandler イベントハンドラ 設定
func (b *CommonServerBuilder) SetEventHandler(event ServerBuildEvents, handler ServerBuildEventHandler) *CommonServerBuilder {
	b.buildEventHandlers[event] = handler
	return b
}

// ClearEventHandler イベントハンドラ クリア
func (b *CommonServerBuilder) ClearEventHandler(event ServerBuildEvents) *CommonServerBuilder {
	delete(b.buildEventHandlers, event)
	return b
}

// GetEventHandler イベントハンドラ 取得
func (b *CommonServerBuilder) GetEventHandler(event ServerBuildEvents) *ServerBuildEventHandler {
	if handler, ok := b.buildEventHandlers[event]; ok {
		return &handler
	}
	return nil
}

// SetDiskEventHandler ディスクイベントハンドラ 設定
func (b *CommonServerBuilder) SetDiskEventHandler(event DiskBuildEvents, handler DiskBuildEventHandler) *CommonServerBuilder {
	b.disk.SetEventHandler(event, handler)
	return b
}

// ClearDiskEventHandler ディスクイベントハンドラ クリア
func (b *CommonServerBuilder) ClearDiskEventHandler(event DiskBuildEvents) *CommonServerBuilder {
	b.disk.ClearEventHandler(event)
	return b
}

// GetDiskEventHandler ディスクイベントハンドラ 取得
func (b *CommonServerBuilder) GetDiskEventHandler(event DiskBuildEvents) *DiskBuildEventHandler {
	return b.disk.GetEventHandler(event)
}

/*---------------------------------------------------------
  for additional disks
---------------------------------------------------------*/

// AddAdditionalDisk 追加ディスク 追加
func (b *CommonServerBuilder) AddAdditionalDisk(diskBuilder *DiskBuilder) *CommonServerBuilder {
	b.additionalDisks = append(b.additionalDisks, diskBuilder)
	return b
}

// ClearAdditionalDisks 追加ディスク クリア
func (b *CommonServerBuilder) ClearAdditionalDisks() *CommonServerBuilder {
	b.additionalDisks = []*DiskBuilder{}
	return b
}

// GetAdditionalDisks 追加ディスク 取得
func (b *CommonServerBuilder) GetAdditionalDisks() []*DiskBuilder {
	return b.additionalDisks
}
