-- Extracts PULUMI_VERSION_MISE and GO_VERSION_MISE from go.mod
-- This replaces the functionality of scripts/get-versions.sh with a Windows-compatible approach

local function read_file(path)
  local file = io.open(path, "r")
  if not file then
    return nil
  end
  local content = file:read("*all")
  file:close()
  return content
end

local function extract_pulumi_version(content, module_path)
  -- Look for lines containing the module path
  for line in content:gmatch("[^\r\n]+") do
    if line:find(module_path, 1, true) then
      -- Extract version starting with 'v' followed by digits
      local version = line:match("v([0-9][^%s]*)")
      if version then
        return version
      end
    end
  end
  return nil
end

local function extract_go_version(content)
  -- Prefer toolchain directive if present
  local toolchain_version = content:match("toolchain%s+go([0-9][^\r\n]*)")
  if toolchain_version then
    return toolchain_version:match("^[^%s]+")
  end

  -- Fall back to go version line
  local go_version = content:match("go%s+([0-9][^\r\n]*)")
  if go_version then
    return go_version:match("^[^%s]+")
  end

  return nil
end

return function(ctx)
  local module_path = "github.com/pulumi/pulumi/pkg/v3"
  local go_mod_path = ctx.options.module_path or ""
  local gomod = "go.mod"

  if go_mod_path ~= "" and go_mod_path ~= "." then
    gomod = go_mod_path .. "/" .. gomod
  end

  local content = read_file(gomod)
  if not content then
    error("missing " .. gomod)
  end

  -- Extract Pulumi version
  local pulumi_version = extract_pulumi_version(content, module_path)
  if not pulumi_version then
    error("failed to determine Pulumi version from " .. gomod)
  end

  -- Extract Go version
  local go_version = extract_go_version(content)
  if not go_version then
    error("failed to determine Go version from " .. gomod)
  end

  return {
    PULUMI_VERSION_MISE = pulumi_version,
    GO_VERSION_MISE = go_version
  }
end
