"use client";

import { useState } from "react";
import { FaChevronDown, FaChevronUp } from "react-icons/fa";

const categories = [
  { id: "all", name: "All Categories" },
  { id: "men", name: "Men's Clothing" },
  { id: "women", name: "Women's Clothing" },
  { id: "accessories", name: "Accessories" },
  { id: "footwear", name: "Footwear" },
];

const colors = [
  { id: "black", name: "Black" },
  { id: "white", name: "White" },
  { id: "gray", name: "Gray" },
  { id: "red", name: "Red" },
  { id: "blue", name: "Blue" },
  { id: "green", name: "Green" },
  { id: "beige", name: "Beige" },
  { id: "brown", name: "Brown" },
];

const sizes = [
  { id: "xs", name: "XS" },
  { id: "s", name: "S" },
  { id: "m", name: "M" },
  { id: "l", name: "L" },
  { id: "xl", name: "XL" },
  { id: "xxl", name: "XXL" },
];

export default function ProductFilters({ filters, onFilterChange }) {
  const [expandedSections, setExpandedSections] = useState({
    categories: true,
    price: true,
    colors: true,
    sizes: true,
  });

  const toggleSection = (section) => {
    setExpandedSections((prev) => ({
      ...prev,
      [section]: !prev[section],
    }));
  };

  const handleCategoryChange = (categoryId) => {
    onFilterChange({ category: categoryId });
  };

  const handlePriceChange = (e) => {
    const { name, value } = e.target;
    onFilterChange({
      price: {
        ...filters.price,
        [name]: Number(value),
      },
    });
  };

  const handleColorChange = (colorId) => {
    const newColors = filters.colors.includes(colorId)
      ? filters.colors.filter((id) => id !== colorId)
      : [...filters.colors, colorId];
    
    onFilterChange({ colors: newColors });
  };

  const handleSizeChange = (sizeId) => {
    const newSizes = filters.sizes.includes(sizeId)
      ? filters.sizes.filter((id) => id !== sizeId)
      : [...filters.sizes, sizeId];
    
    onFilterChange({ sizes: newSizes });
  };

  const handleClearFilters = () => {
    onFilterChange({
      category: "all",
      price: {
        min: 0,
        max: 1000,
      },
      colors: [],
      sizes: [],
    });
  };

  return (
    <div className="bg-white p-4 rounded-lg shadow-sm border">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-lg font-semibold">Filters</h2>
        <button
          className="text-sm text-gray-500 hover:text-black"
          onClick={handleClearFilters}
        >
          Clear all
        </button>
      </div>

      {/* Categories */}
      <div className="mb-6">
        <div
          className="flex justify-between items-center mb-3 cursor-pointer"
          onClick={() => toggleSection("categories")}
        >
          <h3 className="font-medium">Categories</h3>
          {expandedSections.categories ? <FaChevronUp size={14} /> : <FaChevronDown size={14} />}
        </div>
        
        {expandedSections.categories && (
          <div className="space-y-2">
            {categories.map((category) => (
              <div key={category.id} className="flex items-center">
                <input
                  type="radio"
                  id={`category-${category.id}`}
                  name="category"
                  checked={filters.category === category.id}
                  onChange={() => handleCategoryChange(category.id)}
                  className="h-4 w-4 text-black focus:ring-black border-gray-300"
                />
                <label
                  htmlFor={`category-${category.id}`}
                  className="ml-2 text-sm text-gray-700"
                >
                  {category.name}
                </label>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Price Range */}
      <div className="mb-6">
        <div
          className="flex justify-between items-center mb-3 cursor-pointer"
          onClick={() => toggleSection("price")}
        >
          <h3 className="font-medium">Price Range</h3>
          {expandedSections.price ? <FaChevronUp size={14} /> : <FaChevronDown size={14} />}
        </div>
        
        {expandedSections.price && (
          <div className="space-y-4">
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-500">${filters.price.min}</span>
              <span className="text-sm text-gray-500">${filters.price.max}</span>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label htmlFor="min" className="block text-sm mb-1">
                  Min
                </label>
                <input
                  type="number"
                  id="min"
                  name="min"
                  min="0"
                  max={filters.price.max}
                  value={filters.price.min}
                  onChange={handlePriceChange}
                  className="w-full border border-gray-300 rounded p-2 text-sm"
                />
              </div>
              <div>
                <label htmlFor="max" className="block text-sm mb-1">
                  Max
                </label>
                <input
                  type="number"
                  id="max"
                  name="max"
                  min={filters.price.min}
                  value={filters.price.max}
                  onChange={handlePriceChange}
                  className="w-full border border-gray-300 rounded p-2 text-sm"
                />
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Colors */}
      <div className="mb-6">
        <div
          className="flex justify-between items-center mb-3 cursor-pointer"
          onClick={() => toggleSection("colors")}
        >
          <h3 className="font-medium">Colors</h3>
          {expandedSections.colors ? <FaChevronUp size={14} /> : <FaChevronDown size={14} />}
        </div>
        
        {expandedSections.colors && (
          <div className="grid grid-cols-3 gap-2">
            {colors.map((color) => (
              <div
                key={color.id}
                onClick={() => handleColorChange(color.id)}
                className={`p-2 flex flex-col items-center cursor-pointer rounded transition-all ${
                  filters.colors.includes(color.id)
                    ? "bg-gray-100 ring-1 ring-black"
                    : "hover:bg-gray-50"
                }`}
              >
                <div
                  className="w-6 h-6 rounded-full border"
                  style={{ backgroundColor: color.id }}
                ></div>
                <span className="text-xs mt-1">{color.name}</span>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Sizes */}
      <div className="mb-6">
        <div
          className="flex justify-between items-center mb-3 cursor-pointer"
          onClick={() => toggleSection("sizes")}
        >
          <h3 className="font-medium">Sizes</h3>
          {expandedSections.sizes ? <FaChevronUp size={14} /> : <FaChevronDown size={14} />}
        </div>
        
        {expandedSections.sizes && (
          <div className="grid grid-cols-3 gap-2">
            {sizes.map((size) => (
              <div
                key={size.id}
                onClick={() => handleSizeChange(size.id)}
                className={`p-2 flex items-center justify-center cursor-pointer rounded border ${
                  filters.sizes.includes(size.id)
                    ? "bg-black text-white"
                    : "hover:bg-gray-50"
                }`}
              >
                {size.name}
              </div>
            ))}
          </div>
        )}
      </div>

      <button
        className="bg-black text-white w-full py-2 rounded hover:bg-gray-800 transition-colors"
        onClick={() => window.scrollTo({ top: 0, behavior: "smooth" })}
      >
        Apply Filters
      </button>
    </div>
  );
} 