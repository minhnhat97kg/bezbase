# BezBase Frontend

A modern React frontend application featuring responsive design, comprehensive user management, and advanced RBAC interface built with TailwindCSS.

## 🚀 Features

- **Framework**: React 18 with modern hooks and context
- **Styling**: TailwindCSS for responsive, utility-first design
- **Routing**: React Router v6 for navigation
- **Authentication**: JWT token management with automatic refresh
- **Authorization**: RBAC interface with role and permission management
- **State Management**: React Context for global state
- **UI Components**: Custom reusable components with icons
- **API Integration**: Axios-based API client with interceptors
- **Theme Support**: Light/dark theme switching
- **Responsive**: Mobile-first responsive design

## 📁 Project Structure

```
frontend/
├── public/                 # Static assets
│   └── index.html         # HTML template
├── src/                   # Source code
│   ├── components/        # React components
│   │   ├── common/        # Reusable UI components
│   │   │   ├── Icons.js   # Icon components (Lineicons, FontAwesome)
│   │   │   └── Table.js   # Reusable table component
│   │   ├── rbac/          # RBAC-specific components
│   │   │   ├── ActionsList.js        # Actions management
│   │   │   ├── PermissionManager.js  # Permission assignment
│   │   │   ├── ResourcesList.js      # Resources management
│   │   │   ├── RoleForm.js          # Role creation/editing
│   │   │   ├── RolesList.js         # Roles listing
│   │   │   └── UserRoleAssignment.js # User-role assignment
│   │   ├── Header.js      # Application header
│   │   └── Sidebar.js     # Navigation sidebar
│   ├── context/           # React context providers
│   │   ├── AuthContext.js # Authentication state management
│   │   └── ThemeContext.js # Theme management
│   ├── hooks/             # Custom React hooks
│   │   ├── useAuth.js     # Authentication hook
│   │   ├── usePageTitle.js # Page title management
│   │   └── useResourceActionOptions.js # RBAC options hook
│   ├── pages/             # Page components
│   │   ├── Dashboard.js   # Main dashboard
│   │   ├── Login.js      # User login
│   │   ├── Register.js   # User registration
│   │   ├── Profile.js    # User profile management
│   │   ├── RoleManagement.js    # Role management interface
│   │   └── UserManagement.js    # User management interface
│   ├── services/          # API and external services
│   │   └── api.js        # API client with interceptors
│   ├── App.js            # Main application component
│   ├── index.js          # Application entry point
│   └── index.css         # Global styles and Tailwind imports
├── .env.example          # Environment variables template
├── package.json          # Dependencies and scripts
├── tailwind.config.js    # TailwindCSS configuration
├── postcss.config.js     # PostCSS configuration
└── Dockerfile           # Production container image
```

## 🛠️ Quick Start

### Prerequisites

- Node.js 18+ and npm
- Backend API running (see [../backend/README.md](../backend/README.md))

### 1. Local Development Setup

```bash
# Navigate to frontend directory
cd bezbase/frontend

# Copy environment file
cp .env.example .env

# Edit .env with your API configuration
nano .env

# Install dependencies
npm install

# Start development server
npm start
```

### 2. Docker Development

```bash
# Start with Docker Compose (from project root)
docker-compose up frontend

# Or build and run individually
docker build -t bezbase-frontend .
docker run -p 3000:3000 bezbase-frontend
```

### 3. Access the Application

- **Frontend**: http://localhost:3000
- **Auto-reload**: Changes are automatically reflected

## 🔧 Configuration

### Environment Variables

Create a `.env` file in the frontend directory:

```bash
# API Configuration
REACT_APP_API_URL=http://localhost:8080
REACT_APP_ENV=development

# Optional: Feature flags
REACT_APP_ENABLE_THEME_SWITCHING=true
REACT_APP_ENABLE_REGISTRATION=true
```

### Available Scripts

```bash
# Development
npm start          # Start development server with hot reload
npm run build      # Build for production
npm test           # Run test suite
npm run eject      # Eject from Create React App (caution!)

# Code Quality
npm run lint       # Run ESLint (if configured)
npm run format     # Format code with Prettier (if configured)
```

## 🎨 UI Components and Design

### Design System

The application uses a consistent design system built with TailwindCSS:

#### Color Palette
- **Primary**: Blue shades for main actions and navigation
- **Secondary**: Gray shades for text and borders
- **Success**: Green for positive actions
- **Warning**: Yellow for warnings
- **Error**: Red for destructive actions

#### Typography
- **Headings**: `text-xl`, `text-2xl`, `text-3xl` with `font-semibold`
- **Body**: `text-sm`, `text-base` with regular weight
- **Captions**: `text-xs` with `text-gray-500`

#### Spacing
- **Margins**: `m-2`, `m-4`, `m-6` for consistent spacing
- **Padding**: `p-2`, `p-4`, `p-6` for internal spacing
- **Gaps**: `gap-2`, `gap-4`, `gap-6` for flex/grid layouts

### Reusable Components

#### Table Component (`components/common/Table.js`)
```jsx
<Table
  columns={[
    { key: 'name', label: 'Name', sortable: true },
    { key: 'email', label: 'Email' },
    { key: 'actions', label: 'Actions', render: (row) => <Actions /> }
  ]}
  data={users}
  pagination={{ page, pageSize, total }}
  onSort={handleSort}
  onPageChange={handlePageChange}
/>
```

#### Icons Component (`components/common/Icons.js`)
```jsx
// Usage
<UserIcon className="w-5 h-5" />
<EditIcon className="w-4 h-4 text-blue-500" />
<DeleteIcon className="w-4 h-4 text-red-500" />
```

### RBAC Components

#### Role Management
- **RolesList**: Display roles with pagination and search
- **RoleForm**: Create/edit role with validation
- **PermissionManager**: Assign permissions to roles

#### User Management
- **UserManagement**: Complete user CRUD interface
- **UserRoleAssignment**: Assign/remove roles from users

#### Resource & Action Management
- **ResourcesList**: Display available resources
- **ActionsList**: Display available actions

## 🔐 Authentication & Authorization

### Authentication Flow

```jsx
// AuthContext usage
const { user, login, logout, isLoading } = useAuth();

// Login
const handleLogin = async (credentials) => {
  try {
    await login(credentials);
    navigate('/dashboard');
  } catch (error) {
    setError(error.message);
  }
};

// Protected route checking
if (isLoading) return <LoadingSpinner />;
if (!user) return <Navigate to="/login" />;
```

### JWT Token Management

The API client automatically handles JWT tokens:

```javascript
// Auto-attach tokens to requests
axios.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Auto-refresh on token expiry
axios.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Redirect to login
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);
```

### Permission-Based UI

```jsx
// Check permissions before rendering components
const { hasPermission } = useAuth();

return (
  <div>
    {hasPermission('users', 'create') && (
      <button onClick={createUser}>Create User</button>
    )}
    {hasPermission('users', 'read') && (
      <UsersList />
    )}
  </div>
);
```

## 🎯 Key Features

### Dashboard
- Overview cards with statistics
- Recent activity feed
- Quick action buttons
- Responsive grid layout

### User Management
- Complete CRUD operations
- Advanced search and filtering
- Bulk operations
- User status management
- Role assignment interface

### Role Management
- Role CRUD with rich metadata
- Permission assignment
- System role protection
- Active/inactive status management

### Profile Management
- Personal information editing
- Password change
- Theme preferences
- Activity history

### Responsive Design
- Mobile-first approach
- Collapsible sidebar on mobile
- Touch-friendly interactions
- Optimized for tablets and phones

## 🛡️ Security Features

### Input Validation
```jsx
// Form validation example
const [errors, setErrors] = useState({});

const validateForm = (data) => {
  const newErrors = {};
  
  if (!data.email?.match(/^[^\s@]+@[^\s@]+\.[^\s@]+$/)) {
    newErrors.email = 'Invalid email format';
  }
  
  if (data.password?.length < 8) {
    newErrors.password = 'Password must be at least 8 characters';
  }
  
  setErrors(newErrors);
  return Object.keys(newErrors).length === 0;
};
```

### XSS Protection
- Sanitized user inputs
- Content Security Policy headers
- Safe HTML rendering with proper escaping

### CSRF Protection
- API tokens for state-changing operations
- Same-origin policy enforcement

## 🎨 Theming

### Theme Context
```jsx
// ThemeContext usage
const { theme, toggleTheme } = useContext(ThemeContext);

// Apply theme classes
<div className={`${theme === 'dark' ? 'bg-gray-900' : 'bg-white'}`}>
  <button onClick={toggleTheme}>
    {theme === 'dark' ? <SunIcon /> : <MoonIcon />}
  </button>
</div>
```

### TailwindCSS Dark Mode
```javascript
// tailwind.config.js
module.exports = {
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#eff6ff',
          500: '#3b82f6',
          900: '#1e3a8a',
        }
      }
    }
  }
}
```

## 🧪 Testing

### Running Tests
```bash
# Run all tests
npm test

# Run tests in watch mode
npm test -- --watch

# Run tests with coverage
npm test -- --coverage

# Run specific test file
npm test UserManagement.test.js
```

### Testing Structure
```javascript
// Example test
import { render, screen, fireEvent } from '@testing-library/react';
import { AuthProvider } from '../context/AuthContext';
import UserManagement from '../pages/UserManagement';

test('renders user management page', () => {
  render(
    <AuthProvider>
      <UserManagement />
    </AuthProvider>
  );
  
  expect(screen.getByText('User Management')).toBeInTheDocument();
});
```

## 🚀 Deployment

### Production Build

```bash
# Create production build
npm run build

# The build folder contains optimized static files
# Deploy contents to your web server or CDN
```

### Docker Production

```bash
# Build production image
docker build -t bezbase-frontend .

# Run production container
docker run -p 80:80 bezbase-frontend
```

### Environment-Specific Builds

**Development:**
```bash
REACT_APP_API_URL=http://localhost:8080
REACT_APP_ENV=development
```

**Staging:**
```bash
REACT_APP_API_URL=https://api-staging.bezbase.com
REACT_APP_ENV=staging
```

**Production:**
```bash
REACT_APP_API_URL=https://api.bezbase.com
REACT_APP_ENV=production
```

## 📱 Browser Support

- Chrome 88+
- Firefox 85+
- Safari 14+
- Edge 88+
- iOS Safari 14+
- Android Chrome 88+

## 🔧 Development Tools

### Recommended Extensions (VS Code)
- ES7+ React/Redux/React-Native snippets
- Tailwind CSS IntelliSense
- Prettier - Code formatter
- ESLint
- Auto Rename Tag

### Performance Optimization
- Code splitting with React.lazy()
- Image optimization
- Bundle size monitoring
- Service worker for caching (if needed)

### Build Optimization
```javascript
// Analyzing bundle size
npm install -g webpack-bundle-analyzer
npm run build
npx webpack-bundle-analyzer build/static/js/*.js
```

## 🐛 Troubleshooting

### Common Issues

1. **API Connection Failed**
   ```bash
   # Check if backend is running
   curl http://localhost:8080/api/health
   
   # Verify REACT_APP_API_URL in .env
   cat .env | grep REACT_APP_API_URL
   ```

2. **Authentication Not Working**
   ```javascript
   // Check localStorage for token
   console.log(localStorage.getItem('token'));
   
   // Check if token is valid
   // Use browser dev tools → Application → Local Storage
   ```

3. **Styling Issues**
   ```bash
   # Rebuild Tailwind styles
   npm run build:css
   
   # Check if PostCSS is processing correctly
   # Verify tailwind.config.js and postcss.config.js
   ```

4. **Routes Not Working**
   ```javascript
   // Check React Router configuration
   // Ensure BrowserRouter wraps your app
   // Verify route paths match your navigation
   ```

### Debugging

```javascript
// Enable React DevTools
// Add to .env.local for development
REACT_APP_DEBUG=true

// Console debugging
console.log('API Response:', response.data);
console.log('User permissions:', user.permissions);
```

## 📖 Additional Resources

- [React Documentation](https://reactjs.org/docs/)
- [TailwindCSS Documentation](https://tailwindcss.com/docs)
- [React Router Documentation](https://reactrouter.com/)
- [Backend API Documentation](http://localhost:8080/swagger/)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Follow the existing code style and patterns
4. Add tests for new components
5. Test responsiveness on different screen sizes
6. Update documentation if needed
7. Submit a pull request

### Code Style Guidelines

- Use functional components with hooks
- Follow React best practices
- Use meaningful component and variable names
- Keep components small and focused
- Use TypeScript-style prop documentation
- Follow TailwindCSS utility-first approach

---

For questions or support, please refer to the main [project README](../README.md) or create an issue in the repository.