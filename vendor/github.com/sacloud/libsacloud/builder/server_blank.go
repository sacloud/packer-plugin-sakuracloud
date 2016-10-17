package builder

import "github.com/sacloud/libsacloud/sacloud"

// BlankDiskServerBuilder ブランクディスクを利用して構築を行うサーバービルダー
//
// 空のディスクを持ちます。基本的なディスク設定やディスクの追加に対応していますが、ディスクの修正機能には非対応です。
type BlankDiskServerBuilder struct {
	*serverBuilder
}

/*---------------------------------------------------------
  common properties
---------------------------------------------------------*/

// GetServerName サーバー名 取得
func (b *BlankDiskServerBuilder) GetServerName() string {
	return b.serverName
}

// SetServerName サーバー名 設定
func (b *BlankDiskServerBuilder) SetServerName(serverName string) *BlankDiskServerBuilder {
	b.serverName = serverName
	return b
}

// GetCore CPUコア数 取得
func (b *BlankDiskServerBuilder) GetCore() int {
	return b.core
}

// SetCore CPUコア数 設定
func (b *BlankDiskServerBuilder) SetCore(core int) *BlankDiskServerBuilder {
	b.core = core
	return b
}

// GetMemory メモリサイズ(GB単位) 取得
func (b *BlankDiskServerBuilder) GetMemory() int {
	return b.memory
}

// SetMemory メモリサイズ(GB単位) 設定
func (b *BlankDiskServerBuilder) SetMemory(memory int) *BlankDiskServerBuilder {
	b.memory = memory
	return b
}

// IsUseVirtIONetPCI NIC準仮装化モード(virtio)利用フラグ 取得
func (b *BlankDiskServerBuilder) IsUseVirtIONetPCI() bool {
	return b.useVirtIONetPCI
}

// SetUseVirtIONetPCI NIC準仮装化モード(virtio)利用フラグ 設定
func (b *BlankDiskServerBuilder) SetUseVirtIONetPCI(useVirtIONetPCI bool) *BlankDiskServerBuilder {
	b.useVirtIONetPCI = useVirtIONetPCI
	return b
}

// GetDescription 説明 取得
func (b *BlankDiskServerBuilder) GetDescription() string {
	return b.description
}

// SetDescription 説明 設定
func (b *BlankDiskServerBuilder) SetDescription(description string) *BlankDiskServerBuilder {
	b.description = description
	return b
}

// GetIconID アイコンID 取得
func (b *BlankDiskServerBuilder) GetIconID() int64 {
	return b.iconID
}

// SetIconID アイコンID 設定
func (b *BlankDiskServerBuilder) SetIconID(iconID int64) *BlankDiskServerBuilder {
	b.iconID = iconID
	return b
}

// IsBootAfterCreate サーバー作成後すぐに起動フラグ 取得
func (b *BlankDiskServerBuilder) IsBootAfterCreate() bool {
	return b.bootAfterCreate
}

// SetBootAfterCreate サーバー作成後すぐに起動フラグ 設定
func (b *BlankDiskServerBuilder) SetBootAfterCreate(bootAfterCreate bool) *BlankDiskServerBuilder {
	b.bootAfterCreate = bootAfterCreate
	return b
}

// GetTags タグ 取得
func (b *BlankDiskServerBuilder) GetTags() []string {
	return b.Tags
}

// SetTags タグ 設定
func (b *BlankDiskServerBuilder) SetTags(tags []string) *BlankDiskServerBuilder {
	b.Tags = tags
	return b
}

/*---------------------------------------------------------
  for nic functioms
---------------------------------------------------------*/

// ClearNICConnections NIC接続設定 クリア
func (b *BlankDiskServerBuilder) ClearNICConnections() *BlankDiskServerBuilder {
	b.nicConnections = nil
	return b
}

// AddPublicNWConnectedNIC 共有セグメントへの接続追加(注:共有セグメントはeth0のみ接続可能)
func (b *BlankDiskServerBuilder) AddPublicNWConnectedNIC() *BlankDiskServerBuilder {
	b.nicConnections = append(b.nicConnections, "shared")
	return b
}

// AddExistsSwitchConnectedNIC スイッチ or ルーター+スイッチへの接続追加(注:ルーター+スイッチはeth0のみ接続可能)
func (b *BlankDiskServerBuilder) AddExistsSwitchConnectedNIC(switchID string) *BlankDiskServerBuilder {
	b.nicConnections = append(b.nicConnections, switchID)
	return b
}

// AddDisconnectedNIC 切断されたNIC追加
func (b *BlankDiskServerBuilder) AddDisconnectedNIC() *BlankDiskServerBuilder {
	b.nicConnections = append(b.nicConnections, "")
	return b
}

// GetISOImageID ISOイメージ(CDROM)ID 取得
func (b *BlankDiskServerBuilder) GetISOImageID() int64 {
	return b.isoImageID
}

// SetISOImageID ISOイメージ(CDROM)ID 設定
func (b *BlankDiskServerBuilder) SetISOImageID(id int64) *BlankDiskServerBuilder {
	b.isoImageID = id
	return b
}

/*---------------------------------------------------------
  for disk properties
---------------------------------------------------------*/

// GetDiskSize ディスクサイズ(GB単位) 取得
func (b *BlankDiskServerBuilder) GetDiskSize() int {
	return b.disk.GetSize()
}

// SetDiskSize ディスクサイズ(GB単位) 設定
func (b *BlankDiskServerBuilder) SetDiskSize(diskSize int) *BlankDiskServerBuilder {
	b.disk.SetSize(diskSize)
	return b
}

// GetDistantFrom ストレージ隔離対象ディスク 取得
func (b *BlankDiskServerBuilder) GetDistantFrom() []int64 {
	return b.disk.GetDistantFrom()
}

// SetDistantFrom ストレージ隔離対象ディスク 設定
func (b *BlankDiskServerBuilder) SetDistantFrom(distantFrom []int64) *BlankDiskServerBuilder {
	b.disk.SetDistantFrom(distantFrom)
	return b
}

// AddDistantFrom ストレージ隔離対象ディスク 追加
func (b *BlankDiskServerBuilder) AddDistantFrom(diskID int64) *BlankDiskServerBuilder {
	b.disk.AddDistantFrom(diskID)
	return b
}

// ClearDistantFrom ストレージ隔離対象ディスク クリア
func (b *BlankDiskServerBuilder) ClearDistantFrom() *BlankDiskServerBuilder {
	b.disk.ClearDistantFrom()
	return b
}

// GetDiskPlanID ディスクプラン(SSD/HDD) 取得
func (b *BlankDiskServerBuilder) GetDiskPlanID() sacloud.DiskPlanID {
	return b.disk.GetPlanID()
}

// SetDiskPlanID ディスクプラン(SSD/HDD) 設定
func (b *BlankDiskServerBuilder) SetDiskPlanID(diskPlanID sacloud.DiskPlanID) *BlankDiskServerBuilder {
	b.disk.SetPlanID(diskPlanID)
	return b
}

// GetDiskConnection ディスク接続方法(VirtIO/IDE) 取得
func (b *BlankDiskServerBuilder) GetDiskConnection() sacloud.EDiskConnection {
	return b.disk.GetConnection()
}

// SetDiskConnection ディスク接続方法(VirtIO/IDE) 設定
func (b *BlankDiskServerBuilder) SetDiskConnection(diskConnection sacloud.EDiskConnection) *BlankDiskServerBuilder {
	b.disk.SetConnection(diskConnection)
	return b
}

/*---------------------------------------------------------
  for event handler
---------------------------------------------------------*/

// SetEventHandler イベントハンドラ 設定
func (b *BlankDiskServerBuilder) SetEventHandler(event ServerBuildEvents, handler ServerBuildEventHandler) *BlankDiskServerBuilder {
	b.buildEventHandlers[event] = handler
	return b
}

// ClearEventHandler イベントハンドラ クリア
func (b *BlankDiskServerBuilder) ClearEventHandler(event ServerBuildEvents) *BlankDiskServerBuilder {
	delete(b.buildEventHandlers, event)
	return b
}

// GetEventHandler イベントハンドラ 取得
func (b *BlankDiskServerBuilder) GetEventHandler(event ServerBuildEvents) *ServerBuildEventHandler {
	if handler, ok := b.buildEventHandlers[event]; ok {
		return &handler
	}
	return nil
}

// SetDiskEventHandler ディスクイベントハンドラ 設定
func (b *BlankDiskServerBuilder) SetDiskEventHandler(event DiskBuildEvents, handler DiskBuildEventHandler) *BlankDiskServerBuilder {
	b.disk.SetEventHandler(event, handler)
	return b
}

// ClearDiskEventHandler ディスクイベントハンドラ クリア
func (b *BlankDiskServerBuilder) ClearDiskEventHandler(event DiskBuildEvents) *BlankDiskServerBuilder {
	b.disk.ClearEventHandler(event)
	return b
}

// GetDiskEventHandler ディスクイベントハンドラ 取得
func (b *BlankDiskServerBuilder) GetDiskEventHandler(event DiskBuildEvents) *DiskBuildEventHandler {
	return b.disk.GetEventHandler(event)
}

/*---------------------------------------------------------
  for additional disks
---------------------------------------------------------*/

// AddAdditionalDisk 追加ディスク 追加
func (b *BlankDiskServerBuilder) AddAdditionalDisk(diskBuilder *DiskBuilder) *BlankDiskServerBuilder {
	b.additionalDisks = append(b.additionalDisks, diskBuilder)
	return b
}

// ClearAdditionalDisks 追加ディスク クリア
func (b *BlankDiskServerBuilder) ClearAdditionalDisks() *BlankDiskServerBuilder {
	b.additionalDisks = []*DiskBuilder{}
	return b
}

// GetAdditionalDisks 追加ディスク 取得
func (b *BlankDiskServerBuilder) GetAdditionalDisks() []*DiskBuilder {
	return b.additionalDisks
}
