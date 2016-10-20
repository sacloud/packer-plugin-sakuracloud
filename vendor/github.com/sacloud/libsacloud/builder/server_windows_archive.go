package builder

import "github.com/sacloud/libsacloud/sacloud"

// PublicArchiveWindowsServerBuilder Windows系パブリックアーカイブを利用して構築を行うサーバービルダー
//
// 基本的なディスク設定やディスクの追加に対応していますが、ディスクの修正機能には非対応です。
type PublicArchiveWindowsServerBuilder struct {
	*serverBuilder
}

/*---------------------------------------------------------
  common properties
---------------------------------------------------------*/

// GetServerName サーバー名 取得
func (b *PublicArchiveWindowsServerBuilder) GetServerName() string {
	return b.serverName
}

// SetServerName サーバー名 設定
func (b *PublicArchiveWindowsServerBuilder) SetServerName(serverName string) *PublicArchiveWindowsServerBuilder {
	b.serverName = serverName
	return b
}

// GetCore CPUコア数 取得
func (b *PublicArchiveWindowsServerBuilder) GetCore() int {
	return b.core
}

// SetCore CPUコア数 設定
func (b *PublicArchiveWindowsServerBuilder) SetCore(core int) *PublicArchiveWindowsServerBuilder {
	b.core = core
	return b
}

// GetMemory メモリサイズ(GB単位) 取得
func (b *PublicArchiveWindowsServerBuilder) GetMemory() int {
	return b.memory
}

// SetMemory メモリサイズ(GB単位) 設定
func (b *PublicArchiveWindowsServerBuilder) SetMemory(memory int) *PublicArchiveWindowsServerBuilder {
	b.memory = memory
	return b
}

// IsUseVirtIONetPCI NIC準仮装化モード(virtio)利用フラグ 取得
func (b *PublicArchiveWindowsServerBuilder) IsUseVirtIONetPCI() bool {
	return b.useVirtIONetPCI
}

// SetUseVirtIONetPCI NIC準仮装化モード(virtio)利用フラグ 設定
func (b *PublicArchiveWindowsServerBuilder) SetUseVirtIONetPCI(useVirtIONetPCI bool) *PublicArchiveWindowsServerBuilder {
	b.useVirtIONetPCI = useVirtIONetPCI
	return b
}

// GetDescription 説明 取得
func (b *PublicArchiveWindowsServerBuilder) GetDescription() string {
	return b.description
}

// SetDescription 説明 設定
func (b *PublicArchiveWindowsServerBuilder) SetDescription(description string) *PublicArchiveWindowsServerBuilder {
	b.description = description
	return b
}

// GetIconID アイコンID 取得
func (b *PublicArchiveWindowsServerBuilder) GetIconID() int64 {
	return b.iconID
}

// SetIconID アイコンID 設定
func (b *PublicArchiveWindowsServerBuilder) SetIconID(iconID int64) *PublicArchiveWindowsServerBuilder {
	b.iconID = iconID
	return b
}

// IsBootAfterCreate サーバー作成後すぐに起動フラグ 取得
func (b *PublicArchiveWindowsServerBuilder) IsBootAfterCreate() bool {
	return b.bootAfterCreate
}

// SetBootAfterCreate サーバー作成後すぐに起動フラグ 設定
func (b *PublicArchiveWindowsServerBuilder) SetBootAfterCreate(bootAfterCreate bool) *PublicArchiveWindowsServerBuilder {
	b.bootAfterCreate = bootAfterCreate
	return b
}

// GetTags タグ 取得
func (b *PublicArchiveWindowsServerBuilder) GetTags() []string {
	return b.Tags
}

// SetTags タグ 設定
func (b *PublicArchiveWindowsServerBuilder) SetTags(tags []string) *PublicArchiveWindowsServerBuilder {
	b.Tags = tags
	return b
}

/*---------------------------------------------------------
  for nic functioms
---------------------------------------------------------*/

// ClearNICConnections NIC接続設定 クリア
func (b *PublicArchiveWindowsServerBuilder) ClearNICConnections() *PublicArchiveWindowsServerBuilder {
	b.nicConnections = nil
	return b
}

// AddPublicNWConnectedNIC 共有セグメントへの接続追加(注:共有セグメントはeth0のみ接続可能)
func (b *PublicArchiveWindowsServerBuilder) AddPublicNWConnectedNIC() *PublicArchiveWindowsServerBuilder {
	b.nicConnections = append(b.nicConnections, "shared")
	return b
}

// AddExistsSwitchConnectedNIC スイッチ or ルーター+スイッチへの接続追加(注:ルーター+スイッチはeth0のみ接続可能)
func (b *PublicArchiveWindowsServerBuilder) AddExistsSwitchConnectedNIC(switchID string) *PublicArchiveWindowsServerBuilder {
	b.nicConnections = append(b.nicConnections, switchID)
	return b
}

// AddDisconnectedNIC 切断されたNIC追加
func (b *PublicArchiveWindowsServerBuilder) AddDisconnectedNIC() *PublicArchiveWindowsServerBuilder {
	b.nicConnections = append(b.nicConnections, "")
	return b
}

// GetISOImageID ISOイメージ(CDROM)ID 取得
func (b *PublicArchiveWindowsServerBuilder) GetISOImageID() int64 {
	return b.isoImageID
}

// SetISOImageID ISOイメージ(CDROM)ID 設定
func (b *PublicArchiveWindowsServerBuilder) SetISOImageID(id int64) *PublicArchiveWindowsServerBuilder {
	b.isoImageID = id
	return b
}

/*---------------------------------------------------------
  for disk properties
---------------------------------------------------------*/

// GetDiskSize ディスクサイズ(GB単位) 取得
func (b *PublicArchiveWindowsServerBuilder) GetDiskSize() int {
	return b.disk.GetSize()
}

// SetDiskSize ディスクサイズ(GB単位) 設定
func (b *PublicArchiveWindowsServerBuilder) SetDiskSize(diskSize int) *PublicArchiveWindowsServerBuilder {
	b.disk.SetSize(diskSize)
	return b
}

// GetDistantFrom ストレージ隔離対象ディスク 取得
func (b *PublicArchiveWindowsServerBuilder) GetDistantFrom() []int64 {
	return b.disk.GetDistantFrom()
}

// SetDistantFrom ストレージ隔離対象ディスク 設定
func (b *PublicArchiveWindowsServerBuilder) SetDistantFrom(distantFrom []int64) *PublicArchiveWindowsServerBuilder {
	b.disk.SetDistantFrom(distantFrom)
	return b
}

// AddDistantFrom ストレージ隔離対象ディスク 追加
func (b *PublicArchiveWindowsServerBuilder) AddDistantFrom(diskID int64) *PublicArchiveWindowsServerBuilder {
	b.disk.AddDistantFrom(diskID)
	return b
}

// ClearDistantFrom ストレージ隔離対象ディスク クリア
func (b *PublicArchiveWindowsServerBuilder) ClearDistantFrom() *PublicArchiveWindowsServerBuilder {
	b.disk.ClearDistantFrom()
	return b
}

// GetDiskPlanID ディスクプラン(SSD/HDD) 取得
func (b *PublicArchiveWindowsServerBuilder) GetDiskPlanID() sacloud.DiskPlanID {
	return b.disk.GetPlanID()
}

// SetDiskPlanID ディスクプラン(SSD/HDD) 設定
func (b *PublicArchiveWindowsServerBuilder) SetDiskPlanID(diskPlanID sacloud.DiskPlanID) *PublicArchiveWindowsServerBuilder {
	b.disk.SetPlanID(diskPlanID)
	return b
}

// GetDiskConnection ディスク接続方法(VirtIO/IDE) 取得
func (b *PublicArchiveWindowsServerBuilder) GetDiskConnection() sacloud.EDiskConnection {
	return b.disk.GetConnection()
}

// SetDiskConnection ディスク接続方法(VirtIO/IDE) 設定
func (b *PublicArchiveWindowsServerBuilder) SetDiskConnection(diskConnection sacloud.EDiskConnection) *PublicArchiveWindowsServerBuilder {
	b.disk.SetConnection(diskConnection)
	return b
}

/*---------------------------------------------------------
  for disk edit properties
---------------------------------------------------------*/

// GetSourceArchiveID ソースアーカイブID 取得
func (b *PublicArchiveWindowsServerBuilder) GetSourceArchiveID() int64 {
	return b.disk.GetSourceArchiveID()
}

/*---------------------------------------------------------
  for event handler
---------------------------------------------------------*/

// SetEventHandler イベントハンドラ 設定
func (b *PublicArchiveWindowsServerBuilder) SetEventHandler(event ServerBuildEvents, handler ServerBuildEventHandler) *PublicArchiveWindowsServerBuilder {
	b.buildEventHandlers[event] = handler
	return b
}

// ClearEventHandler イベントハンドラ クリア
func (b *PublicArchiveWindowsServerBuilder) ClearEventHandler(event ServerBuildEvents) *PublicArchiveWindowsServerBuilder {
	delete(b.buildEventHandlers, event)
	return b
}

// GetEventHandler イベントハンドラ 取得
func (b *PublicArchiveWindowsServerBuilder) GetEventHandler(event ServerBuildEvents) *ServerBuildEventHandler {
	if handler, ok := b.buildEventHandlers[event]; ok {
		return &handler
	}
	return nil
}

// SetDiskEventHandler ディスクイベントハンドラ 設定
func (b *PublicArchiveWindowsServerBuilder) SetDiskEventHandler(event DiskBuildEvents, handler DiskBuildEventHandler) *PublicArchiveWindowsServerBuilder {
	b.disk.SetEventHandler(event, handler)
	return b
}

// ClearDiskEventHandler ディスクイベントハンドラ クリア
func (b *PublicArchiveWindowsServerBuilder) ClearDiskEventHandler(event DiskBuildEvents) *PublicArchiveWindowsServerBuilder {
	b.disk.ClearEventHandler(event)
	return b
}

// GetDiskEventHandler ディスクイベントハンドラ 取得
func (b *PublicArchiveWindowsServerBuilder) GetDiskEventHandler(event DiskBuildEvents) *DiskBuildEventHandler {
	return b.disk.GetEventHandler(event)
}

/*---------------------------------------------------------
  for additional disks
---------------------------------------------------------*/

// AddAdditionalDisk 追加ディスク 追加
func (b *PublicArchiveWindowsServerBuilder) AddAdditionalDisk(diskBuilder *DiskBuilder) *PublicArchiveWindowsServerBuilder {
	b.additionalDisks = append(b.additionalDisks, diskBuilder)
	return b
}

// ClearAdditionalDisks 追加ディスク クリア
func (b *PublicArchiveWindowsServerBuilder) ClearAdditionalDisks() *PublicArchiveWindowsServerBuilder {
	b.additionalDisks = []*DiskBuilder{}
	return b
}

// GetAdditionalDisks 追加ディスク 取得
func (b *PublicArchiveWindowsServerBuilder) GetAdditionalDisks() []*DiskBuilder {
	return b.additionalDisks
}
