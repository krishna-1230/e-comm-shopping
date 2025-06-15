"use client";

import Link from "next/link";
import { useState, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { useRouter } from "next/navigation";
import { FaShoppingCart, FaHeart, FaUser, FaBars, FaTimes, FaSearch, FaSignOutAlt } from "react-icons/fa";
import { fetchCart } from "@/redux/slices/cartSlice";
import { fetchWishlist } from "@/redux/slices/wishlistSlice";
import { logout, fetchCurrentUser } from "@/redux/slices/authSlice";
import { isAuthenticated } from "@/services/auth-helper";
import toast from "react-hot-toast";

export default function Header() {
  const dispatch = useDispatch();
  const router = useRouter();
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  
  // Get data from Redux store
  const { isAuthenticated: isLoggedIn, user } = useSelector(state => state.auth);
  const { totalItems: cartItems } = useSelector(state => state.cart);
  const { items: wishlistItems } = useSelector(state => state.wishlist);

  useEffect(() => {
    // If user is authenticated, fetch cart and wishlist
    if (isAuthenticated()) {
      dispatch(fetchCurrentUser());
      dispatch(fetchCart());
      dispatch(fetchWishlist());
    }
  }, [dispatch, isLoggedIn]);

  const toggleMenu = () => {
    setIsMenuOpen(!isMenuOpen);
  };

  const handleSearch = (e) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      router.push(`/products?search=${encodeURIComponent(searchQuery)}`);
    }
  };
  
  const handleLogout = () => {
    dispatch(logout());
    toast.success("Logged out successfully");
    router.push('/');
    
    // Close mobile menu if open
    if (isMenuOpen) {
      setIsMenuOpen(false);
    }
  };

  return (
    <header className="bg-white shadow-sm sticky top-0 z-50">
      <div className="container mx-auto px-4 py-4">
        <div className="flex items-center justify-between">
          {/* Logo */}
          <Link href="/" className="text-2xl font-bold">
            StyleSpace
          </Link>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex space-x-8">
            <Link href="/products" className="hover:text-gray-500 transition-colors">
              All Products
            </Link>
            <Link href="/products/men" className="hover:text-gray-500 transition-colors">
              Men
            </Link>
            <Link href="/products/women" className="hover:text-gray-500 transition-colors">
              Women
            </Link>
            <Link href="/products/accessories" className="hover:text-gray-500 transition-colors">
              Accessories
            </Link>
          </nav>

          {/* Search Bar */}
          <div className="hidden md:flex relative">
            <form onSubmit={handleSearch} className="flex items-center">
              <input
                type="text"
                placeholder="Search products..."
                className="border rounded-l-md py-2 px-4 w-64 focus:outline-none focus:ring-2 focus:ring-gray-200"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
              <button
                type="submit"
                className="bg-black text-white p-2 rounded-r-md hover:bg-gray-800 transition-colors"
              >
                <FaSearch />
              </button>
            </form>
          </div>

          {/* User Actions */}
          <div className="hidden md:flex items-center space-x-4">
            <Link href="/wishlist" className="p-2 hover:bg-gray-100 rounded-full transition-colors relative">
              <FaHeart className="text-gray-600" />
              {wishlistItems && wishlistItems.length > 0 && (
                <span className="absolute -top-1 -right-1 bg-black text-white text-xs rounded-full h-5 w-5 flex items-center justify-center">
                  {wishlistItems.length}
                </span>
              )}
            </Link>
            <Link href="/cart" className="p-2 hover:bg-gray-100 rounded-full transition-colors relative">
              <FaShoppingCart className="text-gray-600" />
              {cartItems > 0 && (
                <span className="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center">
                  {cartItems}
                </span>
              )}
            </Link>
            
            {isLoggedIn ? (
              <div className="relative group">
                <button className="p-2 hover:bg-gray-100 rounded-full transition-colors flex items-center">
                  <FaUser className="text-gray-600 mr-2" />
                  <span className="text-sm">{user?.name || "Account"}</span>
                </button>
                <div className="absolute right-0 w-48 mt-2 py-2 bg-white rounded-md shadow-lg hidden group-hover:block border">
                  <Link href="/account" className="block px-4 py-2 text-sm hover:bg-gray-100">
                    My Account
                  </Link>
                  <Link href="/account/orders" className="block px-4 py-2 text-sm hover:bg-gray-100">
                    Orders
                  </Link>
                  <Link href="/account/settings" className="block px-4 py-2 text-sm hover:bg-gray-100">
                    Settings
                  </Link>
                  <button 
                    onClick={handleLogout}
                    className="block w-full text-left px-4 py-2 text-sm hover:bg-gray-100 text-red-600"
                  >
                    Sign Out
                  </button>
                </div>
              </div>
            ) : (
              <Link href="/auth/login" className="p-2 hover:bg-gray-100 rounded-full transition-colors flex items-center">
                <FaUser className="text-gray-600 mr-2" />
                <span className="text-sm">Sign In</span>
              </Link>
            )}
          </div>

          {/* Mobile Menu Button */}
          <button
            onClick={toggleMenu}
            className="p-2 md:hidden rounded-md hover:bg-gray-100 transition-colors"
          >
            {isMenuOpen ? <FaTimes /> : <FaBars />}
          </button>
        </div>
      </div>

      {/* Mobile Menu */}
      {isMenuOpen && (
        <div className="md:hidden bg-white border-t">
          <div className="container mx-auto px-4 py-2">
            <form onSubmit={handleSearch} className="flex items-center my-4">
              <input
                type="text"
                placeholder="Search products..."
                className="border rounded-l-md py-2 px-4 flex-grow focus:outline-none focus:ring-2 focus:ring-gray-200"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
              <button
                type="submit"
                className="bg-black text-white p-2 rounded-r-md hover:bg-gray-800 transition-colors"
              >
                <FaSearch />
              </button>
            </form>
            <nav className="flex flex-col space-y-4 py-4">
              <Link href="/products" className="py-2 hover:text-gray-500 transition-colors">
                All Products
              </Link>
              <Link href="/products/men" className="py-2 hover:text-gray-500 transition-colors">
                Men
              </Link>
              <Link href="/products/women" className="py-2 hover:text-gray-500 transition-colors">
                Women
              </Link>
              <Link href="/products/accessories" className="py-2 hover:text-gray-500 transition-colors">
                Accessories
              </Link>
            </nav>
            <div className="flex flex-col py-4 border-t space-y-4">
              <Link href="/wishlist" className="flex items-center space-x-2">
                <FaHeart /> <span>Wishlist ({wishlistItems?.length || 0})</span>
              </Link>
              <Link href="/cart" className="flex items-center space-x-2">
                <FaShoppingCart /> <span>Cart ({cartItems || 0})</span>
              </Link>
              
              {isLoggedIn ? (
                <>
                  <Link href="/account" className="flex items-center space-x-2">
                    <FaUser /> <span>{user?.name || "My Account"}</span>
                  </Link>
                  <button 
                    onClick={handleLogout}
                    className="flex items-center space-x-2 text-red-600"
                  >
                    <FaSignOutAlt /> <span>Sign Out</span>
                  </button>
                </>
              ) : (
                <Link href="/auth/login" className="flex items-center space-x-2">
                  <FaUser /> <span>Sign In</span>
                </Link>
              )}
            </div>
          </div>
        </div>
      )}
    </header>
  );
} 