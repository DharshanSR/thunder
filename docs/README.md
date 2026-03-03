# Thunder Documentation ⚡

Documentation for **WSO2 Thunder** - a modern identity management suite. This documentation covers installation, configuration, development, and contribution guidelines for the Thunder platform.

## 🌍 Available Languages

Thunder documentation is available in multiple languages:

- [English (US)](https://thunder.dev/) - `en-US` (Default)
- [Español](https://thunder.dev/es-ES/) - `es-ES` (Coming Soon)
- [Français](https://thunder.dev/fr-FR/) - `fr-FR` (Coming Soon)
- [简体中文](https://thunder.dev/zh-CN/) - `zh-CN` (Coming Soon)
- [日本語](https://thunder.dev/ja-JP/) - `ja-JP` (Coming Soon)
- [Português (Brasil)](https://thunder.dev/pt-BR/) - `pt-BR` (Coming Soon)
- [Deutsch](https://thunder.dev/de-DE/) - `de-DE` (Coming Soon)

Want to help translate? See [Translation Guide](./content/community/contributing/translations.md).

## Development

### Working with Translations

```bash
# Extract translatable strings for a locale
pnpm run i18n:extract --locale es-ES

# Start dev server with specific locale
pnpm run dev:locale es-ES

# Build all locales
pnpm run i18n:build
```
