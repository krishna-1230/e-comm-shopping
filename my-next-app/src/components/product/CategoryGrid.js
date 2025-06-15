"use client";

import { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import Link from "next/link";
import { fetchCategories } from "@/redux/slices/productSlice";

// Fallback categories if API fails
const fallbackCategories = [
  {
    id: 1,
    name: "Men's Clothing",
    slug: "men",
    image: "/images/placeholder-category.jpg"
  },
  {
    id: 2,
    name: "Women's Clothing",
    slug: "women",
    image: "/images/placeholder-category.jpg"
  },
  {
    id: 3,
    name: "Accessories",
    slug: "accessories",
    image: "/images/placeholder-category.jpg"
  },
  {
    id: 4,
    name: "Footwear",
    slug: "footwear",
    image: "/images/placeholder-category.jpg"
  }
];

export default function CategoryGrid() {
  const dispatch = useDispatch();
  const { categories, loading, error } = useSelector((state) => state.products);
  const [categoryItems, setCategoryItems] = useState([]);

  useEffect(() => {
    dispatch(fetchCategories());
  }, [dispatch]);

  // Process categories when they change
  useEffect(() => {
    if (categories) {
      // If categories is an array, use it directly
      if (Array.isArray(categories)) {
        setCategoryItems(categories);
      } 
      // If categories is an object with a nested categories array
      else if (categories && categories.categories && Array.isArray(categories.categories)) {
        setCategoryItems(categories.categories);
      }
      // If API response structure is different or failed
      else if (!categories || Object.keys(categories).length === 0) {
        setCategoryItems(fallbackCategories);
      }
    }
  }, [categories]);

  if (loading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {[...Array(4)].map((_, index) => (
          <div key={index} className="animate-pulse">
            <div className="bg-gray-200 h-64 rounded-lg"></div>
          </div>
        ))}
      </div>
    );
  }

  if (error && (!categoryItems || categoryItems.length === 0)) {
    return (
      <div className="text-center py-8">
        <p className="text-red-500">Failed to load categories. Please try again later.</p>
      </div>
    );
  }

  if (!categoryItems || categoryItems.length === 0) {
    // If still no categories available, use fallback
    setCategoryItems(fallbackCategories);
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      {categoryItems.map((category) => (
        <Link
          key={category.id}
          href={`/products/${category.slug || category.id}`}
          className="group relative h-64 rounded-lg overflow-hidden shadow-lg hover:shadow-xl transition-shadow"
        >
          {/* Fallback solid background if image is not available */}
          <div className="absolute inset-0 bg-gray-400"></div>
          
          {/* Category image */}
          <div className="absolute inset-0 bg-cover bg-center group-hover:scale-105 transition-transform duration-500 ease-in-out"
               style={{
                 backgroundImage: `url(${category.image_url || category.image || '/images/placeholder-category.jpg'})`,
                 backgroundSize: 'cover',
                 backgroundPosition: 'center',
               }}>
          </div>
          
          {/* Overlay */}
          <div className="absolute inset-0 bg-black opacity-30 group-hover:opacity-20 transition-opacity"></div>
          
          {/* Content */}
          <div className="absolute inset-0 flex flex-col justify-end p-6">
            <h3 className="text-xl font-bold text-white mb-2">{category.name}</h3>
            <span className="inline-block bg-white bg-opacity-90 text-black text-sm font-medium px-4 py-2 rounded-md opacity-0 group-hover:opacity-100 transform translate-y-2 group-hover:translate-y-0 transition-all">
              Shop Now
            </span>
          </div>
        </Link>
      ))}
    </div>
  );
} 