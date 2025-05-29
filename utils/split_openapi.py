import yaml
import os

def resolve_refs(obj, components):
    if isinstance(obj, dict):
        if "$ref" in obj:
            ref_path = obj["$ref"]
            if ref_path.startswith("#/components/schemas/"):
                schema_name = ref_path.split("/")[-1]
                if schema_name in components["schemas"]:
                    resolved_schema = components["schemas"][schema_name]
                    # 再帰的に参照を解決
                    return resolve_refs(resolved_schema, components)
            elif ref_path.startswith("#/components/parameters/"):
                param_name = ref_path.split("/")[-1]
                if param_name in components["parameters"]:
                    resolved_param = components["parameters"][param_name]
                    return resolve_refs(resolved_param, components)
        for key, value in obj.items():
            obj[key] = resolve_refs(value, components)
    elif isinstance(obj, list):
        for i, item in enumerate(obj):
            obj[i] = resolve_refs(item, components)
    return obj

def split_openapi_by_tags(input_file, output_dir):
    with open(input_file, 'r', encoding='utf-8') as f:
        openapi_spec = yaml.safe_load(f)

    base_info = {
        "openapi": openapi_spec.get("openapi"),
        "info": openapi_spec.get("info"),
        "servers": openapi_spec.get("servers"),
        "security": openapi_spec.get("security"),
    }

    all_tags = openapi_spec.get("tags", [])
    components = openapi_spec.get("components", {})

    # タグごとにパスを分類
    paths_by_tag = {}
    for path, methods in openapi_spec.get("paths", {}).items():
        for method, details in methods.items():
            tags = details.get("tags", [])
            for tag in tags:
                if tag not in paths_by_tag:
                    paths_by_tag[tag] = {}
                if path not in paths_by_tag[tag]:
                    paths_by_tag[tag][path] = {}
                paths_by_tag[tag][path][method] = details

    # 各タグのファイルを作成
    for tag_name, tag_paths in paths_by_tag.items():
        output_spec = base_info.copy()
        output_spec["tags"] = [t for t in all_tags if t["name"] == tag_name]
        output_spec["paths"] = {}

        for path, methods in tag_paths.items():
            output_spec["paths"][path] = {}
            for method, details in methods.items():
                # $refを解決して展開
                resolved_details = resolve_refs(details, components)
                output_spec["paths"][path][method] = resolved_details

        output_file = os.path.join(output_dir, f"{tag_name}.yml")
        with open(output_file, 'w', encoding='utf-8') as f:
            yaml.dump(output_spec, f, allow_unicode=True, sort_keys=False)
        print(f"Created {output_file}")

if __name__ == "__main__":
    input_yaml_file = "authapi.yml"
    output_directory = "auth_yaml_divided"
    split_openapi_by_tags(input_yaml_file, output_directory)