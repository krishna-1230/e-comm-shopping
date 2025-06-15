"use client";

import { useState, useEffect } from "react";
import { useSearchParams } from "next/navigation";
import ProductCard from "@/components/product/ProductCard";
import ProductFilters from "@/components/product/ProductFilters";
import ProductSort from "@/components/product/ProductSort";

// Sample product data (in a real app, this would come from an API)
import { dummyProducts } from "@/components/product/FeaturedProducts";

export default function ProductsPage() {
  const searchParams = useSearchParams();
  const categoryParam = searchParams.get("category");
  
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [filters, setFilters] = useState({
    category: categoryParam || "all",
    price: {
      min: 0,
      max: 1000,
    },
    colors: [],
    sizes: [],
  });
  const [sort, setSort] = useState("newest");

  useEffect(() => {
    // In a real application, you would fetch filtered products from your API
    // Example: const fetchProducts = async () => { const res = await fetch(`/api/products?category=${filters.category}`); ... }
    
    setLoading(true);
    
    // Simulate API delay
    const timer = setTimeout(() => {
      // Filter and sort products based on the current filters and sort
      let filteredProducts = [...dummyProducts];
      
      // Apply category filter
      if (filters.category !== "all") {
        filteredProducts = filteredProducts.filter(
          (product) => product.category === filters.category
        );
      }
      
      // Apply price filter
      filteredProducts = filteredProducts.filter(
        (product) =>
          product.price >= filters.price.min && product.price <= filters.price.max
      );
      
      // Apply color filter
      if (filters.colors.length > 0) {
        filteredProducts = filteredProducts.filter((product) =>
          product.colors.some((color) => filters.colors.includes(color))
        );
      }
      
      // Apply size filter
      if (filters.sizes.length > 0) {
        filteredProducts = filteredProducts.filter((product) =>
          product.sizes.some((size) => filters.sizes.includes(size))
        );
      }
      
      // Apply sorting
      switch (sort) {
        case "price-low":
          filteredProducts.sort((a, b) => a.price - b.price);
          break;
        case "price-high":
          filteredProducts.sort((a, b) => b.price - a.price);
          break;
        case "oldest":
          filteredProducts.sort((a, b) => a.id - b.id);
          break;
        case "newest":
        default:
          filteredProducts.sort((a, b) => b.id - a.id);
          break;
      }
      
      setProducts(filteredProducts);
      setLoading(false);
    }, 500);
    
    return () => clearTimeout(timer);
  }, [filters, sort]);

  const handleFilterChange = (newFilters) => {
    setFilters({ ...filters, ...newFilters });
  };

  const handleSortChange = (newSort) => {
    setSort(newSort);
  };

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Page Title */}
      <h1 className="text-3xl font-bold mb-8">
        {filters.category === "all" ? "All Products" : `${filters.category.charAt(0).toUpperCase() + filters.category.slice(1)} Products`}
      </h1>

      <div className="flex flex-col lg:flex-row gap-8">
        {/* Filters */}
        <div className="lg:w-1/4">
          <ProductFilters filters={filters} onFilterChange={handleFilterChange} />
        </div>

        {/* Products */}
        <div className="lg:w-3/4">
          <div className="mb-6">
            <ProductSort sort={sort} onSortChange={handleSortChange} />
          </div>

          {loading ? (
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
              {[...Array(6)].map((_, index) => (
                <div key={index} className="animate-pulse">
                  <div className="bg-gray-200 h-64 rounded-lg mb-3"></div>
                  <div className="bg-gray-200 h-5 w-2/3 rounded mb-2"></div>
                  <div className="bg-gray-200 h-4 w-1/3 rounded"></div>
                </div>
              ))}
            </div>
          ) : products.length > 0 ? (
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
              {products.map((product) => (
                <ProductCard key={product.id} product={product} />
              ))}
            </div>
          ) : (
            <div className="text-center py-12">
              <h3 className="text-xl font-medium mb-2">No products found</h3>
              <p className="text-gray-500">Try adjusting your filters to find products</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
} 