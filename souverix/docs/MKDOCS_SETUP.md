# MkDocs Material Setup Guide

## Overview

Souverix documentation is built with MkDocs Material and deployed to GitHub Pages via GitHub Actions.

## Structure

```
docs/
├── index.md              # Language chooser
├── en/                   # English documentation
│   ├── index.md
│   ├── platform/
│   ├── architecture/
│   ├── components/
│   ├── compliance/
│   ├── testing/
│   └── operations/
├── fr/                   # French documentation (FR-CA)
│   ├── index.md
│   ├── plateforme/
│   ├── architecture/
│   ├── composants/
│   ├── conformite/
│   ├── tests/
│   └── operations/
└── assets/
    └── images/
```

## Local Development

### Install MkDocs Material

```bash
pip install mkdocs-material[imaging]
```

### Serve Locally

```bash
mkdocs serve
```

Visit: http://127.0.0.1:8000

### Build Site

```bash
mkdocs build
```

Output: `site/` directory

## GitHub Pages Deployment

### Setup

1. Go to repository Settings → Pages
2. Source: Select "GitHub Actions"
3. Save

### Automatic Deployment

The `.github/workflows/pages.yml` workflow:
- Builds on push to `main`
- Deploys to GitHub Pages
- Site available at: `https://dasmlab.github.io/ims/`

### Manual Deployment

```bash
git push origin main
```

## Adding Documentation

### English

1. Create file in `docs/en/<section>/<file>.md`
2. Add to `mkdocs.yml` navigation
3. Commit and push

### French

1. Create file in `docs/fr/<section>/<file>.md`
2. Add to `mkdocs.yml` navigation (Français section)
3. Commit and push

## Navigation Structure

Navigation is defined in `mkdocs.yml`:

```yaml
nav:
  - Home: index.md
  - English:
      - Overview: en/index.md
      - Platform: ...
  - Français (QC):
      - Aperçu: fr/index.md
      - Plateforme: ...
```

## Features

- **Bilingual**: English and French (QC)
- **Material Theme**: Modern, responsive design
- **Dark Mode**: Automatic theme switching
- **Search**: Full-text search in both languages
- **Code Highlighting**: Syntax highlighting for code blocks
- **Admonitions**: Callouts, warnings, notes
- **Tabs**: Tabbed content sections

## Troubleshooting

### Build Fails

Check for:
- Missing files referenced in navigation
- Invalid YAML syntax
- Markdown syntax errors

### Pages Not Updating

- Check GitHub Actions workflow status
- Verify Pages source is set to "GitHub Actions"
- Check workflow logs for errors

---

## End of MkDocs Setup Guide
