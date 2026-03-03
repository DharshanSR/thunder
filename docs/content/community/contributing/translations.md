---
title: Translating Thunder Documentation
sidebar_position: 5
---

# Translating Thunder Documentation

Thunder documentation supports multiple languages through Docusaurus internationalization (i18n). This guide explains how to contribute translations.

## Overview

Thunder uses ISO locale codes (`<LANG>-<COUNTRY>`) for all translations:
- English (US): `en-US` (default)
- Spanish (Spain): `es-ES`
- French (France): `fr-FR`
- Chinese (Simplified): `zh-CN`
- Japanese: `ja-JP`
- Portuguese (Brazil): `pt-BR`
- German (Germany): `de-DE`

## Translation Structure

All translations live in the `docs/i18n/` directory:

```
docs/
├── i18n/
│   ├── es-ES/
│   │   ├── docusaurus-plugin-content-docs/
│   │   │   └── current/
│   │   │       └── guides/
│   │   ├── docusaurus-plugin-content-blog/
│   │   └── code.json
│   ├── fr-FR/
│   └── ...
└── content/ (original English content)
```

## How to Contribute Translations

### 1. Extract Translation Files

```bash
cd docs
pnpm run i18n:extract --locale es-ES
```

This generates translation template files for the specified locale.

### 2. Translate Content

Navigate to `i18n/<locale>/docusaurus-plugin-content-docs/current/` and translate markdown files while:
- Keeping frontmatter structure intact
- Preserving code blocks and technical terms
- Maintaining links and references
- Using consistent terminology

### 3. Translate UI Strings

Edit `i18n/<locale>/code.json` to translate UI elements:

```json
{
  "theme.docs.sidebar.collapseButtonTitle": "Collapse sidebar",
  "theme.docs.sidebar.expandButtonTitle": "Expand sidebar"
}
```

### 4. Test Your Translation

```bash
# Start dev server with your locale
pnpm run dev:locale es-ES

# Build all locales
pnpm run i18n:build
```

### 5. Submit Pull Request

- Fork the repository
- Create a branch: `i18n/add-<locale>-translations`
- Commit your translations
- Open a pull request with:
  - Clear description of what was translated
  - Locale being added or updated
  - Percentage of completion

## Translation Guidelines

### Do's
✅ Keep technical terms in English when appropriate (e.g., "OAuth2", "JWT")
✅ Use native date/time formats for your locale
✅ Maintain consistency with existing translations
✅ Test all links after translation
✅ Keep code examples in English or add locale-specific versions

### Don'ts
❌ Don't translate product names (Thunder, WSO2)
❌ Don't modify code in examples unless necessary for clarity
❌ Don't change file structure or frontmatter keys
❌ Don't translate URLs or GitHub links
❌ Don't mix locale codes

## Review Process

1. **Initial Submission**: Community members submit translations via PR
2. **Review**: Native speakers review for accuracy and consistency
3. **Testing**: CI checks for broken links and build errors
4. **Approval**: Maintainers approve after successful review
5. **Indexing**: Once translations reach 80%+ completion, they're enabled for search engine indexing

## Translation Status

Current translation coverage:

| Locale | Language | Progress | Status |
|--------|----------|----------|--------|
| en-US | English | 100% | Complete |
| es-ES | Spanish | 0% | Not Started |
| fr-FR | French | 0% | Not Started |
| zh-CN | Chinese (Simplified) | 0% | Not Started |
| ja-JP | Japanese | 0% | Not Started |
| pt-BR | Portuguese (Brazil) | 0% | Not Started |
| de-DE | German | 0% | Not Started |

## Resources

- [Docusaurus i18n Documentation](https://docusaurus.io/docs/i18n/introduction)
- [ISO 639-1 Language Codes](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes)
- [ISO 3166-1 Country Codes](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2)
- [Thunder Community Discussions](https://github.com/asgardeo/thunder/discussions)

## Getting Help

- Ask questions in [GitHub Discussions](https://github.com/asgardeo/thunder/discussions)
- Join translation coordination threads
- Review existing translations for reference

Thank you for helping make Thunder accessible to global audiences! 🌍⚡
