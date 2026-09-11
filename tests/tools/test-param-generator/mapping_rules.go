package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// MappingRules はカスタムマッピングルールを保持します
type MappingRules struct {
	TypeNameMappings    map[string]string            `json:"typeNameMappings"`    // リソース名.操作 -> 型名
	FieldNameMappings   map[string]map[string]string `json:"fieldNameMappings"`   // 型名 -> (フィールド名 -> JSONタグ)
	ResourceAliases     map[string][]string          `json:"resourceAliases"`     // リソース名の別名
	FieldAliases        map[string][]string          `json:"fieldAliases"`        // フィールド名の別名
	ModuleSpecificRules map[string]*ModuleRules      `json:"moduleSpecificRules"` // モジュール固有ルール
}

// ModuleRules はモジュール固有のマッピングルールを保持します
// 要件: 5.1, 5.2, 5.3, 5.4, 5.5
type ModuleRules struct {
	TypeNameMappings  map[string]string            `json:"typeNameMappings"`  // リソース名.操作 -> 型名
	FieldNameMappings map[string]map[string]string `json:"fieldNameMappings"` // 型名 -> (フィールド名 -> JSONタグ)
	ResourceAliases   map[string][]string          `json:"resourceAliases"`   // リソース名の別名
	FieldAliases      map[string][]string          `json:"fieldAliases"`      // フィールド名の別名
}

// NewMappingRules は新しいMappingRulesを作成します
func NewMappingRules() *MappingRules {
	return &MappingRules{
		TypeNameMappings:    make(map[string]string),
		FieldNameMappings:   make(map[string]map[string]string),
		ResourceAliases:     make(map[string][]string),
		FieldAliases:        make(map[string][]string),
		ModuleSpecificRules: make(map[string]*ModuleRules),
	}
}

// NewModuleRules は新しいModuleRulesを作成します
func NewModuleRules() *ModuleRules {
	return &ModuleRules{
		TypeNameMappings:  make(map[string]string),
		FieldNameMappings: make(map[string]map[string]string),
		ResourceAliases:   make(map[string][]string),
		FieldAliases:      make(map[string][]string),
	}
}

// LoadMappingRules はマッピングルールをファイルから読み込みます
func LoadMappingRules(filePath string) (*MappingRules, error) {
	// ファイルが存在しない場合は空のルールを返す
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return NewMappingRules(), nil
	}

	// ファイルを読み込む
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read mapping rules file: %w", err)
	}

	// JSONをパース
	var rules MappingRules
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("failed to parse mapping rules JSON: %w", err)
	}

	// nilマップを初期化
	if rules.TypeNameMappings == nil {
		rules.TypeNameMappings = make(map[string]string)
	}
	if rules.FieldNameMappings == nil {
		rules.FieldNameMappings = make(map[string]map[string]string)
	}
	if rules.ResourceAliases == nil {
		rules.ResourceAliases = make(map[string][]string)
	}
	if rules.FieldAliases == nil {
		rules.FieldAliases = make(map[string][]string)
	}
	if rules.ModuleSpecificRules == nil {
		rules.ModuleSpecificRules = make(map[string]*ModuleRules)
	}

	// モジュール固有ルールのnilマップを初期化
	for moduleName, moduleRules := range rules.ModuleSpecificRules {
		if moduleRules == nil {
			rules.ModuleSpecificRules[moduleName] = NewModuleRules()
			continue
		}
		if moduleRules.TypeNameMappings == nil {
			moduleRules.TypeNameMappings = make(map[string]string)
		}
		if moduleRules.FieldNameMappings == nil {
			moduleRules.FieldNameMappings = make(map[string]map[string]string)
		}
		if moduleRules.ResourceAliases == nil {
			moduleRules.ResourceAliases = make(map[string][]string)
		}
		if moduleRules.FieldAliases == nil {
			moduleRules.FieldAliases = make(map[string][]string)
		}
	}

	return &rules, nil
}

// SaveMappingRules はマッピングルールをファイルに保存します
func (mr *MappingRules) SaveMappingRules(filePath string) error {
	// JSONに変換（インデント付き）
	data, err := json.MarshalIndent(mr, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal mapping rules to JSON: %w", err)
	}

	// ファイルに書き込む
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write mapping rules file: %w", err)
	}

	return nil
}

// ApplyTypeNameMapping は型名マッピングを適用します
// resourceName.operation の形式でキーを検索します
// 例: "users.Create" -> "CreateSaasUserJSONRequestBody"
func (mr *MappingRules) ApplyTypeNameMapping(resourceName, operation string) (string, bool) {
	// 直接マッピングを確認
	key := fmt.Sprintf("%s.%s", resourceName, operation)
	if typeName, ok := mr.TypeNameMappings[key]; ok {
		return typeName, true
	}

	// 指定されたリソース名が別名として登録されているか確認
	for originalName, aliases := range mr.ResourceAliases {
		for _, alias := range aliases {
			if alias == resourceName {
				// 別名が見つかった場合、元のリソース名でマッピングを確認
				key := fmt.Sprintf("%s.%s", originalName, operation)
				if typeName, ok := mr.TypeNameMappings[key]; ok {
					return typeName, true
				}
			}
		}
	}

	return "", false
}

// ApplyTypeNameMappingWithModule はモジュール固有ルールを優先して型名マッピングを適用します
// 要件: 5.1, 5.2, 5.3, 5.4, 5.5
func (mr *MappingRules) ApplyTypeNameMappingWithModule(moduleName, resourceName, operation string) (string, bool) {
	// 優先度1: モジュール固有ルール
	if moduleRules, ok := mr.ModuleSpecificRules[moduleName]; ok {
		key := fmt.Sprintf("%s.%s", resourceName, operation)
		if typeName, ok := moduleRules.TypeNameMappings[key]; ok {
			return typeName, true
		}

		// モジュール固有の別名も確認
		for originalName, aliases := range moduleRules.ResourceAliases {
			for _, alias := range aliases {
				if alias == resourceName {
					key := fmt.Sprintf("%s.%s", originalName, operation)
					if typeName, ok := moduleRules.TypeNameMappings[key]; ok {
						return typeName, true
					}
				}
			}
		}
	}

	// 優先度2: グローバルルール
	return mr.ApplyTypeNameMapping(resourceName, operation)
}

// ApplyFieldNameMapping はフィールド名マッピングを適用します
// 指定された型名とフィールド名に対するマッピングを返します
func (mr *MappingRules) ApplyFieldNameMapping(typeName, fieldName string) (string, bool) {
	// 型名に対するフィールドマッピングを確認
	if fieldMappings, ok := mr.FieldNameMappings[typeName]; ok {
		if mappedName, ok := fieldMappings[fieldName]; ok {
			return mappedName, true
		}
	}

	// フィールド名の別名を確認
	if aliases, ok := mr.FieldAliases[fieldName]; ok {
		if fieldMappings, ok := mr.FieldNameMappings[typeName]; ok {
			for _, alias := range aliases {
				if mappedName, ok := fieldMappings[alias]; ok {
					return mappedName, true
				}
			}
		}
	}

	return "", false
}

// ApplyFieldNameMappingWithModule はモジュール固有ルールを優先してフィールド名マッピングを適用します
// 要件: 5.1, 5.2, 5.3, 5.4, 5.5
func (mr *MappingRules) ApplyFieldNameMappingWithModule(moduleName, typeName, fieldName string) (string, bool) {
	// 優先度1: モジュール固有ルール
	if moduleRules, ok := mr.ModuleSpecificRules[moduleName]; ok {
		if fieldMappings, ok := moduleRules.FieldNameMappings[typeName]; ok {
			if mappedName, ok := fieldMappings[fieldName]; ok {
				return mappedName, true
			}
		}

		// モジュール固有の別名も確認
		if aliases, ok := moduleRules.FieldAliases[fieldName]; ok {
			if fieldMappings, ok := moduleRules.FieldNameMappings[typeName]; ok {
				for _, alias := range aliases {
					if mappedName, ok := fieldMappings[alias]; ok {
						return mappedName, true
					}
				}
			}
		}
	}

	// 優先度2: グローバルルール
	return mr.ApplyFieldNameMapping(typeName, fieldName)
}

// GetResourceAliases は指定されたリソース名の別名リストを返します
func (mr *MappingRules) GetResourceAliases(resourceName string) []string {
	if aliases, ok := mr.ResourceAliases[resourceName]; ok {
		return aliases
	}
	return nil
}

// GetFieldAliases は指定されたフィールド名の別名リストを返します
func (mr *MappingRules) GetFieldAliases(fieldName string) []string {
	if aliases, ok := mr.FieldAliases[fieldName]; ok {
		return aliases
	}
	return nil
}

// AddTypeNameMapping は型名マッピングを追加します
func (mr *MappingRules) AddTypeNameMapping(resourceName, operation, typeName string) {
	key := fmt.Sprintf("%s.%s", resourceName, operation)
	mr.TypeNameMappings[key] = typeName
}

// AddFieldNameMapping はフィールド名マッピングを追加します
func (mr *MappingRules) AddFieldNameMapping(typeName, fieldName, mappedName string) {
	if mr.FieldNameMappings[typeName] == nil {
		mr.FieldNameMappings[typeName] = make(map[string]string)
	}
	mr.FieldNameMappings[typeName][fieldName] = mappedName
}

// AddResourceAlias はリソース名の別名を追加します
func (mr *MappingRules) AddResourceAlias(resourceName, alias string) {
	if mr.ResourceAliases[resourceName] == nil {
		mr.ResourceAliases[resourceName] = []string{}
	}
	// 重複チェック
	for _, existing := range mr.ResourceAliases[resourceName] {
		if existing == alias {
			return
		}
	}
	mr.ResourceAliases[resourceName] = append(mr.ResourceAliases[resourceName], alias)
}

// AddFieldAlias はフィールド名の別名を追加します
func (mr *MappingRules) AddFieldAlias(fieldName, alias string) {
	if mr.FieldAliases[fieldName] == nil {
		mr.FieldAliases[fieldName] = []string{}
	}
	// 重複チェック
	for _, existing := range mr.FieldAliases[fieldName] {
		if existing == alias {
			return
		}
	}
	mr.FieldAliases[fieldName] = append(mr.FieldAliases[fieldName], alias)
}

// HasTypeNameMapping は指定されたリソース名と操作に対するマッピングが存在するか確認します
func (mr *MappingRules) HasTypeNameMapping(resourceName, operation string) bool {
	key := fmt.Sprintf("%s.%s", resourceName, operation)
	_, ok := mr.TypeNameMappings[key]
	return ok
}

// HasFieldNameMapping は指定された型名とフィールド名に対するマッピングが存在するか確認します
func (mr *MappingRules) HasFieldNameMapping(typeName, fieldName string) bool {
	if fieldMappings, ok := mr.FieldNameMappings[typeName]; ok {
		_, ok := fieldMappings[fieldName]
		return ok
	}
	return false
}

// GetOrCreateModuleRules はモジュール固有ルールを取得または作成します
// 要件: 5.1, 5.2, 5.3, 5.4, 5.5
func (mr *MappingRules) GetOrCreateModuleRules(moduleName string) *ModuleRules {
	if moduleRules, ok := mr.ModuleSpecificRules[moduleName]; ok {
		return moduleRules
	}

	moduleRules := NewModuleRules()
	mr.ModuleSpecificRules[moduleName] = moduleRules
	return moduleRules
}

// AddModuleTypeNameMapping はモジュール固有の型名マッピングを追加します
// 要件: 5.1, 5.2, 5.3, 5.4, 5.5
func (mr *MappingRules) AddModuleTypeNameMapping(moduleName, resourceName, operation, typeName string) {
	moduleRules := mr.GetOrCreateModuleRules(moduleName)
	key := fmt.Sprintf("%s.%s", resourceName, operation)
	moduleRules.TypeNameMappings[key] = typeName
}

// AddModuleFieldNameMapping はモジュール固有のフィールド名マッピングを追加します
// 要件: 5.1, 5.2, 5.3, 5.4, 5.5
func (mr *MappingRules) AddModuleFieldNameMapping(moduleName, typeName, fieldName, mappedName string) {
	moduleRules := mr.GetOrCreateModuleRules(moduleName)
	if moduleRules.FieldNameMappings[typeName] == nil {
		moduleRules.FieldNameMappings[typeName] = make(map[string]string)
	}
	moduleRules.FieldNameMappings[typeName][fieldName] = mappedName
}

// AddModuleResourceAlias はモジュール固有のリソース名別名を追加します
// 要件: 5.1, 5.2, 5.3, 5.4, 5.5
func (mr *MappingRules) AddModuleResourceAlias(moduleName, resourceName, alias string) {
	moduleRules := mr.GetOrCreateModuleRules(moduleName)
	if moduleRules.ResourceAliases[resourceName] == nil {
		moduleRules.ResourceAliases[resourceName] = []string{}
	}
	// 重複チェック
	for _, existing := range moduleRules.ResourceAliases[resourceName] {
		if existing == alias {
			return
		}
	}
	moduleRules.ResourceAliases[resourceName] = append(moduleRules.ResourceAliases[resourceName], alias)
}

// AddModuleFieldAlias はモジュール固有のフィールド名別名を追加します
// 要件: 5.1, 5.2, 5.3, 5.4, 5.5
func (mr *MappingRules) AddModuleFieldAlias(moduleName, fieldName, alias string) {
	moduleRules := mr.GetOrCreateModuleRules(moduleName)
	if moduleRules.FieldAliases[fieldName] == nil {
		moduleRules.FieldAliases[fieldName] = []string{}
	}
	// 重複チェック
	for _, existing := range moduleRules.FieldAliases[fieldName] {
		if existing == alias {
			return
		}
	}
	moduleRules.FieldAliases[fieldName] = append(moduleRules.FieldAliases[fieldName], alias)
}
