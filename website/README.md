# vLLM Semantic Router Documentation

This directory contains the Docusaurus-based documentation website for the vLLM Semantic Router project

## 🚀 Quick Start

### Prerequisites

- Node.js 18+ 
- npm or yarn

### Development

Start the development server with hot reload:

```bash
# From project root
make docs-dev

# Or manually
cd website && npm start
```

The site will be available at http://localhost:3000

### Production Build

Build the static site for production:

```bash
# From project root
make docs-build

# Or manually
cd website && npm run build
```

### Preview Production Build

Serve the production build locally:

```bash
# From project root
make docs-serve

# Or manually
cd website && npm run serve
```

## 🎨 Features

### ✨ Modern Tech-Inspired Design

- **Dark theme by default** with neon blue/green accents
- **Glassmorphism effects** with backdrop blur and transparency
- **Gradient backgrounds** and animated hover effects
- **Responsive design** optimized for all devices

### 🔧 Enhanced Functionality

- **Mermaid diagram support** with dark theme optimization
- **Advanced code highlighting** with multiple language support
- **Interactive navigation** with smooth animations
- **Search functionality** (ready for Algolia integration)

### 📱 User Experience

- **Fast loading** with optimized builds
- **Accessible design** following WCAG guidelines
- **Mobile-first** responsive layout
- **SEO optimized** with proper meta tags

## 📁 Project Structure

```
website/
├── docs/                   # Documentation content (Markdown files)
├── src/
│   ├── components/        # Custom React components
│   ├── css/              # Global styles and theme
│   └── pages/            # Custom pages (homepage, etc.)
├── static/               # Static assets (images, icons, etc.)
├── docusaurus.config.js  # Main configuration
├── sidebars.js          # Navigation structure
└── package.json         # Dependencies and scripts
```

## 🛠️ Customization

### Themes and Colors
Edit `src/css/custom.css` to modify:

- Color scheme and gradients
- Typography and spacing
- Component styling
- Animations and effects

### Navigation
Update `sidebars.js` to modify:

- Documentation structure
- Category organization
- Page ordering

### Site Configuration
Modify `docusaurus.config.js` for:

- Site metadata
- Plugin configuration
- Theme settings
- Build options

## 📚 Available Commands

| Command | Description |
|---------|-------------|
| `make docs-dev` | Start development server |
| `make docs-build` | Build for production |
| `make docs-serve` | Preview production build |
| `make docs-clean` | Clear build cache |

## 🔗 Links

- **Live Preview**: http://localhost:3000 (when running)
- **Docusaurus Docs**: https://docusaurus.io/docs
- **Main Project**: ../README.md
