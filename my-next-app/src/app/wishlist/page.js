"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import Image from "next/image";
import { FaTrash, FaShoppingCart } from "react-icons/fa";
import toast from "react-hot-toast";

// Sample wishlist items (in a real app, this would come from Redux or context)
const initialWishlistItems = [
  {
    id: 3,
    name: "Casual Denim Jacket",
    slug: "casual-denim-jacket",
    price: 79.99,
    discount: 0,
    image: "/images/product-3.jpg",
    color: "Blue",
    size: "M",
  },
  {
    id: 4,
    name: "Leather Crossbody Bag",
    slug: "leather-crossbody-bag",
    price: 89.99,
    discount: 15,
    image: "/images/product-4.jpg",
    color: "Brown",
    size: "",
  },
];

export default function WishlistPage() {
  const [wishlistItems, setWishlistItems] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // In a real application, you would fetch wishlist from your API/Redux store
    const fetchWishlist = async () => {
      setLoading(true);
      
      // Simulating API delay
      await new Promise(resolve => setTimeout(resolve, 500));
      
      setWishlistItems(initialWishlistItems);
      setLoading(false);
    };
    
    fetchWishlist();
  }, []);

  const removeItem = (id) => {
    setWishlistItems(wishlistItems.filter(item => item.id !== id));
    toast.success("Item removed from wishlist");
  };

  const addToCart = (item) => {
    // In a real application, you would dispatch an action to add to cart
    toast.success(`${item.name} added to your cart!`);
  };

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold mb-8">Your Wishlist</h1>
        <div className="animate-pulse">
          <div className="h-20 bg-gray-200 rounded-md mb-4"></div>
          <div className="h-20 bg-gray-200 rounded-md mb-4"></div>
          <div className="h-40 bg-gray-200 rounded-md"></div>
        </div>
      </div>
    );
  }

  if (wishlistItems.length === 0) {
    return (
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold mb-8">Your Wishlist</h1>
        <div className="bg-white p-6 rounded-lg shadow-sm border text-center">
          <p className="mb-4 text-xl">Your wishlist is empty</p>
          <p className="mb-8 text-gray-500">Save items you like by clicking the heart icon on product pages.</p>
          <Link
            href="/products"
            className="bg-black text-white px-6 py-3 rounded-md hover:bg-gray-800 transition-colors"
          >
            Browse Products
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Your Wishlist</h1>
      
      <div className="bg-white rounded-lg shadow-sm border overflow-hidden">
        <div className="grid grid-cols-12 gap-4 p-4 border-b bg-gray-50 hidden md:grid">
          <div className="col-span-7">
            <span className="font-medium">Product</span>
          </div>
          <div className="col-span-2 text-center">
            <span className="font-medium">Price</span>
          </div>
          <div className="col-span-3 text-right">
            <span className="font-medium">Actions</span>
          </div>
        </div>
        
        {/* Wishlist items list */}
        {wishlistItems.map(item => (
          <div key={item.id} className="grid grid-cols-1 md:grid-cols-12 gap-4 p-4 border-b items-center">
            {/* Product info */}
            <div className="col-span-7 flex items-center space-x-4">
              <div className="w-20 h-20 bg-gray-100 rounded flex-shrink-0">
                {/* Replace with proper image when available */}
                <div
                  className="w-full h-full"
                  style={{
                    backgroundImage: `url(${item.image})`,
                    backgroundSize: "cover",
                    backgroundPosition: "center",
                  }}
                />
              </div>
              <div>
                <Link
                  href={`/products/${item.slug}`}
                  className="font-medium hover:text-gray-500 transition-colors"
                >
                  {item.name}
                </Link>
                {(item.color || item.size) && (
                  <p className="text-sm text-gray-500">
                    {item.color && `Color: ${item.color}`}
                    {item.color && item.size && ', '}
                    {item.size && `Size: ${item.size}`}
                  </p>
                )}
              </div>
            </div>
            
            {/* Price */}
            <div className="md:col-span-2 flex justify-between md:justify-center items-center">
              <span className="md:hidden font-medium">Price:</span>
              {item.discount > 0 ? (
                <div className="flex flex-col md:items-center">
                  <span className="text-red-600">
                    ${(item.price * (1 - item.discount / 100)).toFixed(2)}
                  </span>
                  <span className="text-gray-500 text-sm line-through">
                    ${item.price.toFixed(2)}
                  </span>
                </div>
              ) : (
                <span>${item.price.toFixed(2)}</span>
              )}
            </div>
            
            {/* Actions */}
            <div className="md:col-span-3 flex justify-between md:justify-end items-center space-x-2">
              <button
                onClick={() => addToCart(item)}
                className="bg-black text-white px-4 py-2 rounded-md hover:bg-gray-800 transition-colors flex items-center"
              >
                <FaShoppingCart className="mr-2" />
                <span className="hidden sm:inline">Add to Cart</span>
                <span className="sm:hidden">Add</span>
              </button>
              <button
                onClick={() => removeItem(item.id)}
                className="text-gray-400 hover:text-red-500 transition-colors p-2"
                title="Remove from wishlist"
              >
                <FaTrash size={16} />
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
} 