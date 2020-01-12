package md

import (
	"strings"
	"sync"

	"github.com/ggoop/goutils/repositories"
	"github.com/ggoop/goutils/utils"
)

type MDEnum struct {
	EntityID string `gorm:"size:50;primary_key:uix" json:"entity_id"`
	Domain   string `gorm:"size:50" json:"domain"`
	ID       string `gorm:"size:50;primary_key:uix" json:"id"`
	Name     string `gorm:"size:50" json:"name"`
	Sequence int    `json:"sequence"`
}

func (t MDEnum) TableName() string {
	return "md_enums"
}
func (s *MDEnum) MD() *Mder {
	return &Mder{ID: "md.enum", Domain: md_domain, Name: "枚举", Type: TYPE_ENUM}
}

var enumCache map[string]MDEnum

type EnumSv struct {
	repo *repositories.MysqlRepo
	*sync.Mutex
}

/**
* 创建服务实例
 */
func NewEnumSv(repo *repositories.MysqlRepo) *EnumSv {
	return &EnumSv{repo: repo, Mutex: &sync.Mutex{}}
}
func GetEnum(typeId string, values ...string) *MDEnum {
	if enumCache == nil || typeId == "" || values == nil || len(values) == 0 {
		return nil
	}
	for _, v := range values {
		if v, ok := enumCache[strings.ToLower(typeId+":"+v)]; ok {
			return &v
		}
	}
	return nil
}

func (s *EnumSv) InitCache() {
	enumCache = make(map[string]MDEnum)
	items, _ := s.Get()
	for _, v := range items {
		enumCache[strings.ToLower(v.EntityID+":"+v.ID)] = v
		enumCache[strings.ToLower(v.EntityID+":"+v.Name)] = v
	}
}

func (s *EnumSv) Get() ([]MDEnum, error) {
	items := make([]MDEnum, 0)
	if err := s.repo.Model(&MDEnum{}).Where("entity_id in (?)", s.repo.Model(MDEntity{}).Select("id").Where("type=?", "enum").SubQuery()).Order("entity_id").Order("sequence").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
func (s *EnumSv) GetBy(typeId string) ([]MDEnum, error) {
	items := make([]MDEnum, 0)
	if err := s.repo.Model(&MDEnum{}).Where("entity_id=?", typeId).Order("sequence").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
func (s *EnumSv) UpdateOrCreate(enum *MDEnum) (*MDEnum, error) {
	entity := MDEntity{}
	if enum.EntityID == "" {
		return nil, nil
	}
	s.repo.Model(entity).Where("id=?", enum.EntityID).Order("id").Take(&entity)
	if entity.ID == "" {
		entity.ID = enum.EntityID
		entity.Code = enum.EntityID
		entity.Name = enum.EntityID
		entity.Type = TYPE_ENUM
		entity.Domain = enum.Domain
		s.repo.Create(&entity)
	}
	old := MDEnum{}

	if s.repo.Where("entity_id=? and id=?", enum.EntityID, enum.ID).Order("id").Take(&old).RecordNotFound() {
		s.repo.Create(enum)
	} else {
		updates := utils.Map{}
		if old.Name != enum.Name && enum.Name != "" {
			updates["Name"] = enum.Name
		}
		if old.Sequence != enum.Sequence && enum.Sequence > 0 {
			updates["Sequence"] = enum.Sequence
		}
		if len(updates) > 0 {
			s.repo.Model(&old).Updates(updates)
		}
	}
	return enum, nil
}
