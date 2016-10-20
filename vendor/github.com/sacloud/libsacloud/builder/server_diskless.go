package builder

// DisklessServerBuilder ディスクレス サーバービルダー
//
// ディスクレスのサーバーを構築します。 ディスク関連の設定に非対応です。
type DisklessServerBuilder struct {
	*serverBuilder
}

/*---------------------------------------------------------
  common properties
---------------------------------------------------------*/

// GetServerName サーバー名 取得
func (b *DisklessServerBuilder) GetServerName() string {
	return b.serverName
}

// SetServerName サーバー名 設定
func (b *DisklessServerBuilder) SetServerName(serverName string) *DisklessServerBuilder {
	b.serverName = serverName
	return b
}

// GetCore CPUコア数 取得
func (b *DisklessServerBuilder) GetCore() int {
	return b.core
}

// SetCore CPUコア数 設定
func (b *DisklessServerBuilder) SetCore(core int) *DisklessServerBuilder {
	b.core = core
	return b
}

// GetMemory メモリサイズ(GB単位) 取得
func (b *DisklessServerBuilder) GetMemory() int {
	return b.memory
}

// SetMemory メモリサイズ(GB単位) 設定
func (b *DisklessServerBuilder) SetMemory(memory int) *DisklessServerBuilder {
	b.memory = memory
	return b
}

// IsUseVirtIONetPCI NIC準仮装化モード(virtio)利用フラグ 取得
func (b *DisklessServerBuilder) IsUseVirtIONetPCI() bool {
	return b.useVirtIONetPCI
}

// SetUseVirtIONetPCI NIC準仮装化モード(virtio)利用フラグ 設定
func (b *DisklessServerBuilder) SetUseVirtIONetPCI(useVirtIONetPCI bool) *DisklessServerBuilder {
	b.useVirtIONetPCI = useVirtIONetPCI
	return b
}

// GetDescription 説明 取得
func (b *DisklessServerBuilder) GetDescription() string {
	return b.description
}

// SetDescription 説明 設定
func (b *DisklessServerBuilder) SetDescription(description string) *DisklessServerBuilder {
	b.description = description
	return b
}

// GetIconID アイコンID 取得
func (b *DisklessServerBuilder) GetIconID() int64 {
	return b.iconID
}

// SetIconID アイコンID 設定
func (b *DisklessServerBuilder) SetIconID(iconID int64) *DisklessServerBuilder {
	b.iconID = iconID
	return b
}

// IsBootAfterCreate サーバー作成後すぐに起動フラグ 取得
func (b *DisklessServerBuilder) IsBootAfterCreate() bool {
	return b.bootAfterCreate
}

// SetBootAfterCreate サーバー作成後すぐに起動フラグ 設定
func (b *DisklessServerBuilder) SetBootAfterCreate(bootAfterCreate bool) *DisklessServerBuilder {
	b.bootAfterCreate = bootAfterCreate
	return b
}

// GetTags タグ 取得
func (b *DisklessServerBuilder) GetTags() []string {
	return b.Tags
}

// SetTags タグ 設定
func (b *DisklessServerBuilder) SetTags(tags []string) *DisklessServerBuilder {
	b.Tags = tags
	return b
}

/*---------------------------------------------------------
  for nic functioms
---------------------------------------------------------*/

// ClearNICConnections NIC接続設定 クリア
func (b *DisklessServerBuilder) ClearNICConnections() *DisklessServerBuilder {
	b.nicConnections = nil
	return b
}

// AddPublicNWConnectedNIC 共有セグメントへの接続追加(注:共有セグメントはeth0のみ接続可能)
func (b *DisklessServerBuilder) AddPublicNWConnectedNIC() *DisklessServerBuilder {
	b.nicConnections = append(b.nicConnections, "shared")
	return b
}

// AddExistsSwitchConnectedNIC スイッチ or ルーター+スイッチへの接続追加(注:ルーター+スイッチはeth0のみ接続可能)
func (b *DisklessServerBuilder) AddExistsSwitchConnectedNIC(switchID string) *DisklessServerBuilder {
	b.nicConnections = append(b.nicConnections, switchID)
	return b
}

// AddDisconnectedNIC 切断されたNIC追加
func (b *DisklessServerBuilder) AddDisconnectedNIC() *DisklessServerBuilder {
	b.nicConnections = append(b.nicConnections, "")
	return b
}

// GetISOImageID ISOイメージ(CDROM)ID 取得
func (b *DisklessServerBuilder) GetISOImageID() int64 {
	return b.isoImageID
}

// SetISOImageID ISOイメージ(CDROM)ID 設定
func (b *DisklessServerBuilder) SetISOImageID(id int64) *DisklessServerBuilder {
	b.isoImageID = id
	return b
}

/*---------------------------------------------------------
  for event handler
---------------------------------------------------------*/

// SetEventHandler イベントハンドラ 設定
func (b *DisklessServerBuilder) SetEventHandler(event ServerBuildEvents, handler ServerBuildEventHandler) *DisklessServerBuilder {
	b.buildEventHandlers[event] = handler
	return b
}

// ClearEventHandler イベントハンドラ クリア
func (b *DisklessServerBuilder) ClearEventHandler(event ServerBuildEvents) *DisklessServerBuilder {
	delete(b.buildEventHandlers, event)
	return b
}

// GetEventHandler イベントハンドラ 取得
func (b *DisklessServerBuilder) GetEventHandler(event ServerBuildEvents) *ServerBuildEventHandler {
	if handler, ok := b.buildEventHandlers[event]; ok {
		return &handler
	}
	return nil
}
