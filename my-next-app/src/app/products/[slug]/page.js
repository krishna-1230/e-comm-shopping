"use client";

import { useState, useEffect } from "react";
import { useParams } from "next/navigation";
import Image from "next/image";
import Link from "next/link";
import toast from "react-hot-toast";
import { FaArrowLeft, FaHeart, FaRegHeart } from "react-icons/fa";

// Import dummy products data
import { dummyProducts } from "@/components/product/FeaturedProducts";

export default function ProductDetailPage() {
  const { slug } = useParams();
  const [product, setProduct] = useState(null);
  const [loading, setLoading] = useState(true);
  const [selectedColor, setSelectedColor] = useState("");
  const [selectedSize, setSelectedSize] = useState("");
  const [quantity, setQuantity] = useState(1);
  const [isWishlisted, setIsWishlisted] = useState(false);
  const [activeTab, setActiveTab] = useState("description");
  const [relatedProducts, setRelatedProducts] = useState([]);

  useEffect(() => {
    // In a real app, you would fetch the product from your API
    // For now, we'll use the dummy data
    const fetchProduct = async () => {
      setLoading(true);
      
      // Simulate API delay
      await new Promise(resolve => setTimeout(resolve, 500));
      
      const foundProduct = dummyProducts.find(p => p.slug === slug);
      
      if (foundProduct) {
        setProduct(foundProduct);
        setSelectedColor(foundProduct.colors[0] || "");
        setSelectedSize(foundProduct.sizes[0] || "");
        
        // Get related products (same category)
        const related = dummyProducts.filter(
          p => p.category === foundProduct.category && p.id !== foundProduct.id
        ).slice(0, 4);
        
        setRelatedProducts(related);
      }
      
      setLoading(false);
    };
    
    fetchProduct();
  }, [slug]);

  const handleQuantityChange = (e) => {
    const value = parseInt(e.target.value, 10);
    if (value > 0) {
      setQuantity(value);
    }
  };

  const handleDecreaseQuantity = () => {
    if (quantity > 1) {
      setQuantity(quantity - 1);
    }
  };

  const handleIncreaseQuantity = () => {
    setQuantity(quantity + 1);
  };

  const handleColorSelect = (color) => {
    setSelectedColor(color);
  };

  const handleSizeSelect = (size) => {
    setSelectedSize(size);
  };

  const handleToggleWishlist = () => {
    setIsWishlisted(!isWishlisted);
    
    if (!isWishlisted) {
      toast.success(`${product.name} added to your wishlist!`);
    } else {
      toast.success(`${product.name} removed from your wishlist!`);
    }
  };

  const handleAddToCart = () => {
    // In a real application, you would dispatch an action to add to cart
    toast.success(`${product.name} added to your cart!`);
  };

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="animate-pulse">
          <div className="h-6 w-48 bg-gray-200 mb-8"></div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div className="bg-gray-200 h-96 rounded-lg"></div>
            <div className="space-y-4">
              <div className="h-8 bg-gray-200 w-3/4 rounded"></div>
              <div className="h-6 bg-gray-200 w-1/4 rounded"></div>
              <div className="h-4 bg-gray-200 w-full rounded"></div>
              <div className="h-4 bg-gray-200 w-full rounded"></div>
              <div className="h-10 bg-gray-200 w-full rounded mt-8"></div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (!product) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center py-12">
          <h1 className="text-2xl font-bold mb-4">Product Not Found</h1>
          <p className="mb-6">The product you are looking for does not exist or has been removed.</p>
          <Link
            href="/products"
            className="bg-black text-white px-6 py-2 rounded-md hover:bg-gray-800 transition-colors"
          >
            Browse Products
          </Link>
        </div>
      </div>
    );
  }

  const discountedPrice = product.discount > 0
    ? product.price - (product.price * product.discount / 100)
    : product.price;

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Breadcrumb */}
      <div className="flex items-center mb-8 text-sm">
        <Link href="/products" className="flex items-center text-gray-500 hover:text-black">
          <FaArrowLeft className="mr-2" />
          Back to Products
        </Link>
      </div>

      {/* Product Detail */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-16">
        {/* Product Image */}
        <div className="rounded-lg overflow-hidden bg-gray-100">
          {/* Placeholder for product image */}
          <div
            className="h-[500px] w-full"
            style={{
              backgroundImage: `url(${product.images[0] || ''})`,
              backgroundSize: "cover",
              backgroundPosition: "center",
            }}
          />
        </div>

        {/* Product Info */}
        <div>
          <h1 className="text-3xl font-bold mb-2">{product.name}</h1>
          
          {/* Price */}
          <div className="mb-6">
            {product.discount > 0 ? (
              <div className="flex items-center gap-2">
                <span className="text-2xl font-semibold text-red-600">${discountedPrice.toFixed(2)}</span>
                <span className="text-gray-500 line-through">${product.price.toFixed(2)}</span>
                <span className="bg-red-100 text-red-800 text-xs font-medium px-2 py-0.5 rounded">
                  {product.discount}% OFF
                </span>
              </div>
            ) : (
              <span className="text-2xl font-semibold">${product.price.toFixed(2)}</span>
            )}
          </div>
          
          {/* Description */}
          <p className="text-gray-600 mb-8">
            Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, urna eu 
            tincidunt consectetur, nisi nunc ultricies nisi, eget ultricies nisl nisl eu nisl.
          </p>
          
          {/* Colors */}
          {product.colors.length > 0 && (
            <div className="mb-6">
              <h3 className="text-sm font-medium mb-2">Color: {selectedColor}</h3>
              <div className="flex gap-2">
                {product.colors.map(color => (
                  <button
                    key={color}
                    onClick={() => handleColorSelect(color)}
                    className={`w-8 h-8 rounded-full border ${
                      selectedColor === color
                        ? "ring-2 ring-black ring-offset-1"
                        : "hover:scale-110 transition-transform"
                    }`}
                    style={{ backgroundColor: color.toLowerCase() }}
                    title={color}
                  />
                ))}
              </div>
            </div>
          )}
          
          {/* Sizes */}
          {product.sizes.length > 0 && (
            <div className="mb-8">
              <h3 className="text-sm font-medium mb-2">Size: {selectedSize}</h3>
              <div className="flex flex-wrap gap-2">
                {product.sizes.map(size => (
                  <button
                    key={size}
                    onClick={() => handleSizeSelect(size)}
                    className={`py-2 px-3 border rounded-md text-sm ${
                      selectedSize === size
                        ? "border-black bg-black text-white"
                        : "hover:border-gray-400 transition-colors"
                    }`}
                  >
                    {size}
                  </button>
                ))}
              </div>
            </div>
          )}
          
          {/* Quantity */}
          <div className="mb-8">
            <h3 className="text-sm font-medium mb-2">Quantity:</h3>
            <div className="flex items-center w-32">
              <button
                onClick={handleDecreaseQuantity}
                className="bg-gray-100 border border-gray-300 rounded-l-md py-2 px-4 hover:bg-gray-200"
              >
                -
              </button>
              <input
                type="number"
                min="1"
                value={quantity}
                onChange={handleQuantityChange}
                className="w-full text-center border-t border-b border-gray-300 py-2"
              />
              <button
                onClick={handleIncreaseQuantity}
                className="bg-gray-100 border border-gray-300 rounded-r-md py-2 px-4 hover:bg-gray-200"
              >
                +
              </button>
            </div>
          </div>
          
          {/* Actions */}
          <div className="flex gap-4">
            <button
              onClick={handleAddToCart}
              className="flex-grow bg-black text-white py-3 px-6 rounded-md hover:bg-gray-800 transition-colors"
            >
              Add to Cart
            </button>
            <button
              onClick={handleToggleWishlist}
              className="p-3 border rounded-md hover:bg-gray-50 transition-colors"
            >
              {isWishlisted ? (
                <FaHeart className="text-red-500" />
              ) : (
                <FaRegHeart />
              )}
            </button>
          </div>
          
          {/* Product Info Tabs */}
          <div className="mt-12">
            <div className="border-b border-gray-200">
              <nav className="flex -mb-px">
                <button
                  onClick={() => setActiveTab("description")}
                  className={`px-4 py-2 text-sm font-medium ${
                    activeTab === "description"
                      ? "border-b-2 border-black text-black"
                      : "text-gray-500 hover:text-black"
                  }`}
                >
                  Description
                </button>
                <button
                  onClick={() => setActiveTab("details")}
                  className={`px-4 py-2 text-sm font-medium ${
                    activeTab === "details"
                      ? "border-b-2 border-black text-black"
                      : "text-gray-500 hover:text-black"
                  }`}
                >
                  Additional Details
                </button>
                <button
                  onClick={() => setActiveTab("reviews")}
                  className={`px-4 py-2 text-sm font-medium ${
                    activeTab === "reviews"
                      ? "border-b-2 border-black text-black"
                      : "text-gray-500 hover:text-black"
                  }`}
                >
                  Reviews (0)
                </button>
              </nav>
            </div>
            <div className="py-4">
              {activeTab === "description" && (
                <div>
                  <p className="text-gray-600">
                    Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam auctor, 
                    nunc vitae ultricies ultricies, nunc nisl ultricies nunc, vitae ultricies
                    nisl nisl vitae nisl. Nullam auctor, nunc vitae ultricies ultricies, nunc nisl
                    ultricies nunc, vitae ultricies nisl nisl vitae nisl.
                  </p>
                </div>
              )}
              {activeTab === "details" && (
                <div className="space-y-2">
                  <div className="grid grid-cols-2 gap-2">
                    <span className="text-gray-600">Material</span>
                    <span>100% Cotton</span>
                  </div>
                  <div className="grid grid-cols-2 gap-2">
                    <span className="text-gray-600">Care Instructions</span>
                    <span>Machine wash, tumble dry low</span>
                  </div>
                  <div className="grid grid-cols-2 gap-2">
                    <span className="text-gray-600">Origin</span>
                    <span>Imported</span>
                  </div>
                </div>
              )}
              {activeTab === "reviews" && (
                <div className="text-center py-6">
                  <p className="text-gray-500">No reviews yet</p>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Related Products */}
      {relatedProducts.length > 0 && (
        <div className="mt-16">
          <h2 className="text-2xl font-bold mb-6">Related Products</h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
            {relatedProducts.map((product) => (
              <Link key={product.id} href={`/products/${product.slug}`}>
                <div className="product-card rounded-lg overflow-hidden bg-white shadow hover:shadow-md transition-shadow">
                  {/* Product Image */}
                  <div className="relative h-64 bg-gray-100">
                    <div
                      className="h-full w-full"
                      style={{
                        backgroundImage: `url(${product.images[0] || ''})`,
                        backgroundSize: "cover",
                        backgroundPosition: "center",
                      }}
                    />
                  </div>
                  
                  {/* Product Info */}
                  <div className="p-4">
                    <h3 className="text-base font-medium mb-1 truncate-2">{product.name}</h3>
                    <div className="flex items-center gap-2">
                      {product.discount > 0 ? (
                        <>
                          <span className="text-red-600 font-semibold">
                            ${(product.price - (product.price * product.discount / 100)).toFixed(2)}
                          </span>
                          <span className="text-gray-500 text-sm line-through">
                            ${product.price.toFixed(2)}
                          </span>
                        </>
                      ) : (
                        <span className="font-semibold">${product.price.toFixed(2)}</span>
                      )}
                    </div>
                  </div>
                </div>
              </Link>
            ))}
          </div>
        </div>
      )}
    </div>
  );
} 