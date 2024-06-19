package model

import (
	"time"

	"gorm.io/gorm"
)

type ATA struct {
	ID         string       `gorm:"primarykey;type:varchar(20);not null" json:"id"` // 26-01
	Name       string       `gorm:"type:varchar(20);" json:"name"`
	Equipments []*Equipment `gorm:"foreignkey:ATAID" json:"equipments"`
}

func (ATA) TableName() string {
	return "tbl_ata"
}

type Equipment struct {
	ID            uint    `gorm:"primarykey" json:"id"`
	ATAID         string  `gorm:"type:varchar(20);not null" json:"ATA_id"`          // 26-01
	Name          string  `gorm:"type:varchar(30);not null" json:"name"`            // 设备名称
	ConditionText string  `gorm:"type:varchar(912);not null" json:"condition_text"` //前端在选中后需要在界面上显示的
	ProtocolType  uint    `json:"protocol_type"`                                    // 接口类型
	Parts         []*Part `gorm:"many2many:rel_equipment_part" json:"part_list"`
}

func (Equipment) TableName() string {
	return "tbl_equipment"
}

type Part struct {
	ID          string       `gorm:"primarykey" json:"id"` // string?
	Description string       `gorm:"type:varchar(200)" json:"description"`
	Equipments  []*Equipment `gorm:"many2many:rel_equipment_part" json:"equipments"`
}

func (Part) TableName() string {
	return "tbl_part"
}

// 定义 RelEquipmentPart 结构体，表示多对多关系
type RelEquipmentPart struct {
	ID           uint          `gorm:"primarykey" json:"id"`
	EquipmentID  uint          `json:"equipment_id"`
	PartID       string        `json:"part_id"`
	PartLoadLogs []PartLoadLog `gorm:"foreignkey:RelEquipmentPartID;association_foreignkey:ID" json:"part_load_logs"`
}

// 设置 MemberSystem 和 Part 之间的关联
func (RelEquipmentPart) TableName() string {
	return "rel_equipment_part"
}

type BaseModel struct {
	ID        uint `gorm:"primarykey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type PartLoadLog struct {
	BaseModel
	RelEquipmentPartID uint      `gorm:"not null" json:"rel_equipment_part_id"`
	LoadStatus         string    `gorm:"type:varchar(20);not null" json:"load_status"`
	LoadProgress       uint8     `gorm:"type:varchar(20);not null" json:"load_progress"`
	StartTime          time.Time `json:"start_time"`
	ElapsedTime        uint      `json:"elapsed_time"`
	Detail             string    `gorm:"type:varchar(200);not null" json:"detail"`
}

func (PartLoadLog) TableName() string {
	return "tbl_partloadlog"
}

type ATAEvent struct {
	ATA
	ATAList []ATA `json:"ATAList"`
	Total   int64 `json:"total"`
}

// EditUser 编辑用户信息
func (ctl *ATAEvent) GetAll() error {
	// err = db.Find(&ctl.ATAList).Error
	err = db.Model(&ATA{}).Preload("Equipments").Find(&ctl.ATAList).Count(&ctl.Total).Error
	return err
}

type EquipmentEvent struct {
	Equipment
	EquipmentList []Equipment `json:"equipment_list"`
}

type PartEvent struct {
	Part
	PartList []Part `json:"part_list"`
}

func (ctl *EquipmentEvent) GetListByEquipment(id int) error {
	err = db.Preload("Parts").Where("ID = ?", id).First(&ctl.Equipment).Error
	return err
}

type LoadATAEquipmentParam struct {
	EquipmentID uint   `json:"equipment_id"`
	PartlistId  []uint `json:"partlist_id"`
}

type PartLoadLogEvent struct {
	Equipment
	Part Part
	PartLoadLog
	RelEquipmentPart     RelEquipmentPart
	RelEquipmentPartList []RelEquipmentPart
	EquipmentList        []Equipment
	PartList             []Part        `json:"part_list"`
	PartLoadLogList      []PartLoadLog `json:"part_load_log_list"`
}

func (ctl *PartLoadLogEvent) SaveLog(data LoadATAEquipmentParam) error {
	err := db.Model(&RelEquipmentPart{}).Where("equipment_id = ?", data.EquipmentID).Where("part_id in (?)", data.PartlistId).Find(&ctl.RelEquipmentPartList).Error
	for _, v := range ctl.RelEquipmentPartList {
		db.Create(&PartLoadLog{
			RelEquipmentPartID: v.ID,
			LoadStatus:         "Loading",
			LoadProgress:       0,
			StartTime:          time.Now(),
			ElapsedTime:        0,
			Detail:             "Loading",
		})
	}
	return err
}

func (ctl *PartLoadLogEvent) GetAll() error {
	err = db.Model(&PartLoadLog{}).Preload("RelEquipmentPart").Preload("RelEquipmentPart.Equipment").Preload("RelEquipmentPart.Part").Find(&ctl.PartLoadLogList).Error
	return err
}
