# Hugo Cli

1. Create a poster simply.
2. Create a posts easy by config.

## Usage

Hugo Ai Helper cli will read the `.hcli_config.yaml` config.

```bash
hcli gen post -n {Tempalte name} xxx
```

### Generate images

Generate AI-powered images using VolcEngine's image generation API:

```bash
hcli gen pic "a beautiful landscape with mountains and lake"
```

**Configuration requirements:**
- Set `AI.APIKey` in `.hcli_config.yaml` or `VOLC_API_KEY` environment variable
- Set `AI.PicModel` in config or `PIC_MODEL_ID` environment variable
- Optional: Configure custom prompts in `AI.CustomPrompt`

**Example config:**
```yaml
AI:
  Provider: "volc"
  APIKey: "your_volc_api_key"
  ThinkModel: "volc-think-model-id"
  PicModel: "volc-pic-model-id"
  CustomPrompt:
    pic_summary_prompt_key: "Custom image description prompt"
```

### Show all templates

```bash
hcli gen
```

# Install & Make

```bash
# make release from source
make releases

# Make target OS and Arch
# make hcli_{OS}_{ARCH}
make hcli_darwin_arm64


```

