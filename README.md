# StyleSpace E-Commerce Platform

A modern e-commerce platform built with Next.js, Go Fiber, and PostgreSQL.

## Project Structure

```
ecommerce/
├── backend/           # Go Fiber backend API
├── my-next-app/       # Next.js frontend application
└── README.md
```

## Features

- 🛍️ Product browsing and filtering
- 🔍 Advanced search functionality
- 🛒 Shopping cart management
- ❤️ Wishlist functionality
- 👤 User authentication and profiles
- 📦 Order management
- 💳 Secure payment processing
- 📱 Responsive design
- 🌐 RESTful API architecture

## Tech Stack

### Frontend
- Next.js 14
- React
- Tailwind CSS
- Redux Toolkit
- React Hot Toast
- Geist Font

### Backend
- Go Fiber
- PostgreSQL
- GORM
- JWT Authentication
- CORS enabled

## Prerequisites

- Node.js 18+ 
- Go 1.21+
- PostgreSQL 15+
- npm or yarn

## Getting Started

### Backend Setup

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Install Go dependencies:
   ```bash
   go mod download
   ```

3. Create a `.env` file in the backend directory with the following variables:
   ```
   DB_HOST=localhost
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_NAME=ecommerce
   DB_PORT=5432
   JWT_SECRET=your_jwt_secret
   PORT=8080
   ```

4. Run the backend server:
   ```bash
   go run main.go
   ```

### Frontend Setup

1. Navigate to the frontend directory:
   ```bash
   cd my-next-app
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Create a `.env.local` file in the frontend directory:
   ```
   NEXT_PUBLIC_API_URL=http://localhost:8080
   ```

4. Run the development server:
   ```bash
   npm run dev
   ```

The application will be available at:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - User login
- `GET /api/auth/profile` - Get user profile

### Products
- `GET /api/products` - Get all products
- `GET /api/products/:id` - Get product by ID
- `GET /api/products/category/:category` - Get products by category

### Cart
- `GET /api/cart` - Get user's cart
- `POST /api/cart` - Add item to cart
- `PUT /api/cart/:id` - Update cart item
- `DELETE /api/cart/:id` - Remove item from cart

### Orders
- `GET /api/orders` - Get user's orders
- `POST /api/orders` - Create new order
- `GET /api/orders/:id` - Get order details

## Development

### Code Style
- Frontend: ESLint and Prettier
- Backend: Go standard formatting

### Git Workflow
1. Create a new branch for each feature
2. Write meaningful commit messages
3. Create a pull request for review
4. Merge after approval

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support, email support@stylespace.com or create an issue in the repository. "# e-comm-shopping" 
