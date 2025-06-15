"use client";

import { useState, useRef, useEffect } from "react";
import { FaSort, FaCheck } from "react-icons/fa";

const sortOptions = [
  { id: "newest", name: "Newest" },
  { id: "oldest", name: "Oldest" },
  { id: "price-low", name: "Price: Low to High" },
  { id: "price-high", name: "Price: High to Low" },
];

export default function ProductSort({ sort, onSortChange }) {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef(null);
  
  const selectedOption = sortOptions.find((option) => option.id === sort);

  const toggleDropdown = () => {
    setIsOpen(!isOpen);
  };

  const handleSortSelect = (sortId) => {
    onSortChange(sortId);
    setIsOpen(false);
  };

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setIsOpen(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  return (
    <div className="flex justify-end">
      <div className="relative" ref={dropdownRef}>
        <button
          onClick={toggleDropdown}
          className="flex items-center space-x-2 px-4 py-2 border border-gray-300 rounded-md bg-white text-sm hover:bg-gray-50"
        >
          <FaSort />
          <span>Sort by: {selectedOption?.name}</span>
        </button>

        {isOpen && (
          <div className="absolute right-0 z-10 mt-1 w-56 bg-white rounded-md shadow-lg border border-gray-200">
            <div className="py-1">
              {sortOptions.map((option) => (
                <button
                  key={option.id}
                  onClick={() => handleSortSelect(option.id)}
                  className={`flex items-center justify-between w-full px-4 py-2 text-left text-sm hover:bg-gray-100 ${
                    option.id === sort ? "bg-gray-50" : ""
                  }`}
                >
                  <span>{option.name}</span>
                  {option.id === sort && <FaCheck className="text-black" />}
                </button>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
} 