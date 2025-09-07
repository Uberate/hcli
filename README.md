# Hugo Cli

1. Create a poster simply.
2. Create a posts easy by config.

## Usage

Hugo Ai Helper cli will read the `.hcli_config.yaml` config.

```bash
hcli gen post -n {Tempalte name} xxx
```

### Generate images

Generate AI-powered images from template content using VolcEngine's image generation API:

```bash
# Generate image from a markdown file using a specific template
hcli gen pic -n {template_name} filename.md

# Example: Generate image for a blog post
hcli gen pic -n blog_post my-post.md
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
```### Generate Posts

Create new blog posts using AI-powered templates:

```bash
# Generate a new post with a specific template and title
hcli gen post -n blog_template "My New Blog Post"

# Generate a post with custom front matter fields
hcli gen post -n article_template -f "tags=tech,ai;categories=development" "AI in Web Development"

# Generate a post from a YAML config file
hcli gen post -c post_config.yaml
```

**Configuration for Posts:**
- Ensure `AI.ThinkModel` is set in `.hcli_config.yaml` or `THINK_MODEL_ID` environment variable
- Customize post templates in the `Templates.Post` section of the config
- Define default front matter fields in `Defaults.Post`

### Advanced Configuration

The `.hcli_config.yaml` supports additional configuration options:

```yaml
# Example advanced configuration
AI:
  Provider: "volc"  # Supported providers: volc, openai (experimental)
  APIKey: "your_api_key"
  ThinkModel: "volc-think-model-id"
  PicModel: "volc-pic-model-id"
  Timeout: 30  # Request timeout in seconds
Templates:
  Post:
    blog_template:
      prompt: "Write a blog post about {title} with SEO-friendly content"
      front_matter:
        draft: false
        author: "Your Name"
Defaults:
  Post:
    tags: ["ai", "hugo"]
    categories: ["development"]
```

**Environment Variables:**
- `HCLI_CONFIG`: Path to custom config file (overrides default .hcli_config.yaml)
- `DEBUG_MODE`: Set to "true" to enable verbose logging


### Show all templates

```bash
hcli gen
```

# Install & Make

## Quick Install (One-liner)

```bash
# Install with curl (Linux/Mac)
curl -sSL https://raw.githubusercontent.com/Uberate/hcli/refs/heads/main/scripts/install_mac.sh | bash

```

```
# Or download and install manually
curl -LO https://github.com/uberate/hugo-ai-helper/releases/latest/download/hcli_$(uname -s)_$(uname -m).tar.gz
tar -xzf hcli_*.tar.gz
sudo mv hcli /usr/local/bin/
```

```
# Verify installation
hcli --version
```

## Build from Source

```bash
# make release from source
make releases

# Make target OS and Arch
# make hcli_{OS}_{ARCH}
make hcli_darwin_arm64
```

