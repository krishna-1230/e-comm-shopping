"use client";

import { useState, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { fetchProducts } from "@/redux/slices/productSlice";
import ProductCard from "./ProductCard";

// Fallback products if API fails
export const dummyProducts = [
  {
    id: 1,
    name: "Classic White T-Shirt",
    slug: "classic-white-t-shirt",
    price: 24.99,
    discount: 0,
    images: ["/images/product-1.jpg"],
    colors: ["White", "Black", "Gray"],
    sizes: ["S", "M", "L", "XL"],
    isFeatured: true,
    category: "men",
  },
  {
    id: 2,
    name: "Summer Floral Dress",
    slug: "summer-floral-dress",
    price: 49.99,
    discount: 10,
    images: ["/images/product-2.jpg"],
    colors: ["Blue", "Red"],
    sizes: ["S", "M", "L"],
    isFeatured: true,
    category: "women",
  },
  {
    id: 3,
    name: "Casual Denim Jacket",
    slug: "casual-denim-jacket",
    price: 79.99,
    discount: 0,
    images: ["/images/product-3.jpg"],
    colors: ["Blue"],
    sizes: ["S", "M", "L", "XL"],
    isFeatured: true,
    category: "men",
  },
  {
    id: 4,
    name: "Leather Crossbody Bag",
    slug: "leather-crossbody-bag",
    price: 89.99,
    discount: 15,
    images: ["/images/product-4.jpg"],
    colors: ["Brown", "Black"],
    sizes: [],
    isFeatured: true,
    category: "accessories",
  }
];

export default function FeaturedProducts() {
  const dispatch = useDispatch();
  const { products, loading, error } = useSelector((state) => state.products);
  const [featuredProducts, setFeaturedProducts] = useState([]);

  useEffect(() => {
    // Fetch products with isFeatured=true parameter
    dispatch(fetchProducts({ isFeatured: true }));
  }, [dispatch]);

  useEffect(() => {
    // Process products when they change
    if (products) {
      // If products is an array, use it directly
      if (Array.isArray(products) && products.length > 0) {
        // If API returns too many products, limit to first 8
        setFeaturedProducts(products.slice(0, 8));
      } 
      // If products is an object with a nested products array
      else if (products && products.products && Array.isArray(products.products)) {
        setFeaturedProducts(products.products.slice(0, 8));
      }
      // If API call fails or returns empty
      else {
        // Use fallback dummy products
        setFeaturedProducts(dummyProducts);
      }
    }
  }, [products]);

  if (error && featuredProducts.length === 0) {
    // Show error message only if we don't have fallback products
    return (
      <div className="text-center py-8">
        <p className="text-red-500">Failed to load products. Please try again later.</p>
      </div>
    );
  }

  if (loading && featuredProducts.length === 0) {
    return (
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        {[...Array(4)].map((_, index) => (
          <div key={index} className="animate-pulse">
            <div className="bg-gray-200 h-64 rounded-lg mb-3"></div>
            <div className="bg-gray-200 h-5 w-2/3 rounded mb-2"></div>
            <div className="bg-gray-200 h-4 w-1/3 rounded"></div>
          </div>
        ))}
      </div>
    );
  }

  // If no featured products found
  if (featuredProducts.length === 0) {
    // Use fallback dummy products if no featured products are available
    setFeaturedProducts(dummyProducts);
    return (
      <div className="text-center py-8">
        <p className="text-gray-500">Loading featured products...</p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
      {featuredProducts.map((product) => (
        <ProductCard key={product.id} product={product} />
      ))}
    </div>
  );
} 