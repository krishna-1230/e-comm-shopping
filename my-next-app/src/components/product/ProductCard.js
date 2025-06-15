"use client";

import { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import Link from "next/link";
import { FaHeart, FaRegHeart, FaShoppingCart } from "react-icons/fa";
import toast from "react-hot-toast";
import { addItemToCart } from "@/redux/slices/cartSlice";
import { addItemToWishlist, removeWishlistItem, toggleWishlistProduct } from "@/redux/slices/wishlistSlice";
import { isAuthenticated } from "@/services/auth-helper";

export default function ProductCard({ product }) {
  const dispatch = useDispatch();
  const { items: wishlistItems } = useSelector(state => state.wishlist);
  const { loading: cartLoading } = useSelector(state => state.cart);
  const { loading: wishlistLoading } = useSelector(state => state.wishlist);
  const [isInWishlist, setIsInWishlist] = useState(false);
  
  // Ensure product has all necessary properties
  const safeProduct = {
    id: product?.id || 0,
    name: product?.name || 'Product Name',
    slug: product?.slug || product?.id?.toString() || 'product',
    price: product?.price || product?.base_price || 0,
    discount: product?.discount || product?.discount_percentage || 0,
    images: product?.images || [],
    colors: product?.colors || [],
    sizes: product?.sizes || [],
  };
  
  // Get primary image URL or find first available image
  const getImageUrl = () => {
    // Try to find product_images array with image_url property
    if (product?.product_images && Array.isArray(product.product_images) && product.product_images.length > 0) {
      // Try to find primary image first
      const primaryImage = product.product_images.find(img => img.is_primary);
      if (primaryImage && primaryImage.image_url) return primaryImage.image_url;
      
      // Otherwise use first image
      if (product.product_images[0].image_url) return product.product_images[0].image_url;
    }
    
    // Fall back to images array if it exists
    if (Array.isArray(safeProduct.images) && safeProduct.images.length > 0) {
      return safeProduct.images[0];
    }
    
    // Use image property if it exists
    if (product?.image) return product.image;
    
    // Default placeholder
    return '/images/placeholder-product.jpg';
  };
  
  const imageUrl = getImageUrl();
  const discountedPrice = safeProduct.discount > 0 
    ? safeProduct.price - (safeProduct.price * safeProduct.discount / 100)
    : safeProduct.price;

  useEffect(() => {
    // Check if the product is in the wishlist
    if (wishlistItems && Array.isArray(wishlistItems) && wishlistItems.length > 0) {
      setIsInWishlist(wishlistItems.some(item => item.id === safeProduct.id));
    }
  }, [wishlistItems, safeProduct.id]);
  
  const handleAddToCart = (e) => {
    e.preventDefault();
    e.stopPropagation();
    
    // Check if user is authenticated
    if (!isAuthenticated()) {
      toast.error("Please log in to add items to your cart");
      return;
    }
    
    // Create cart item
    const cartItem = {
      productId: safeProduct.id,
      quantity: 1,
      // If the product has colors or sizes, you might want to select default ones
      ...(safeProduct.colors && safeProduct.colors.length > 0 ? { 
        color: Array.isArray(safeProduct.colors) ? safeProduct.colors[0] : 
               (typeof safeProduct.colors === 'object' ? Object.values(safeProduct.colors)[0] : safeProduct.colors) 
      } : {})
    };
    
    dispatch(addItemToCart(cartItem))
      .unwrap()
      .then(() => {
        toast.success(`${safeProduct.name} added to your cart!`);
      })
      .catch((error) => {
        toast.error(error || "Failed to add item to cart");
      });
  };
  
  const handleToggleWishlist = (e) => {
    e.preventDefault();
    e.stopPropagation();
    
    // Check if user is authenticated
    if (!isAuthenticated()) {
      toast.error("Please log in to manage your wishlist");
      return;
    }
    
    if (isInWishlist) {
      dispatch(removeWishlistItem(safeProduct.id))
        .unwrap()
        .then(() => {
          toast.success(`${safeProduct.name} removed from your wishlist!`);
        })
        .catch((error) => {
          toast.error(error || "Failed to remove from wishlist");
        });
    } else {
      dispatch(addItemToWishlist({ productId: safeProduct.id }))
        .unwrap()
        .then(() => {
          toast.success(`${safeProduct.name} added to your wishlist!`);
        })
        .catch((error) => {
          toast.error(error || "Failed to add to wishlist");
        });
    }
  };

  return (
    <Link href={`/products/${safeProduct.slug}`} className="group">
      <div className="product-card relative rounded-lg overflow-hidden bg-white shadow hover:shadow-md transition-shadow">
        {/* Product Image */}
        <div className="relative h-64 bg-gray-100">
          {/* Placeholder for product image */}
          <div
            className="h-full w-full bg-gray-200"
            style={{
              backgroundImage: `url(${imageUrl})`,
              backgroundSize: "cover",
              backgroundPosition: "center",
            }}
          />
          
          {/* Discount badge */}
          {safeProduct.discount > 0 && (
            <div className="absolute top-2 left-2 bg-red-500 text-white text-xs font-bold px-2 py-1 rounded">
              {safeProduct.discount}% OFF
            </div>
          )}
          
          {/* Action buttons */}
          <div className="absolute top-2 right-2">
            <button
              onClick={handleToggleWishlist}
              disabled={wishlistLoading}
              className={`p-2 bg-white rounded-full shadow hover:bg-gray-100 transition-colors ${
                wishlistLoading ? 'opacity-50 cursor-not-allowed' : ''
              }`}
            >
              {isInWishlist ? (
                <FaHeart className="text-red-500" />
              ) : (
                <FaRegHeart />
              )}
            </button>
          </div>
          
          {/* Quick add to cart button */}
          <div className="absolute bottom-0 left-0 right-0 bg-black bg-opacity-0 group-hover:bg-opacity-70 flex justify-center py-2 translate-y-full group-hover:translate-y-0 transition-all duration-300">
            <button
              onClick={handleAddToCart}
              disabled={cartLoading}
              className={`flex items-center space-x-2 text-white hover:text-gray-200 transition-colors ${
                cartLoading ? 'opacity-50 cursor-not-allowed' : ''
              }`}
            >
              <FaShoppingCart />
              <span>Add to Cart</span>
            </button>
          </div>
        </div>
        
        {/* Product Info */}
        <div className="p-4">
          <h3 className="text-base font-medium mb-1 truncate-2">{safeProduct.name}</h3>
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              {safeProduct.discount > 0 ? (
                <>
                  <span className="text-red-600 font-semibold">${discountedPrice.toFixed(2)}</span>
                  <span className="text-gray-500 text-sm line-through">${safeProduct.price.toFixed(2)}</span>
                </>
              ) : (
                <span className="font-semibold">${safeProduct.price.toFixed(2)}</span>
              )}
            </div>
            {/* Color options preview */}
            {safeProduct.colors && safeProduct.colors.length > 0 && (
              <div className="flex items-center gap-1">
                {Array.isArray(safeProduct.colors) ? (
                  safeProduct.colors.slice(0, 3).map((color, index) => (
                    <div
                      key={index}
                      className="w-3 h-3 rounded-full border border-gray-300"
                      style={{ 
                        backgroundColor: typeof color === 'string' ? 
                          (color.startsWith('#') ? color : color.toLowerCase()) : 
                          (color.color_hex || '#CCCCCC')
                      }}
                    />
                  ))
                ) : (
                  <div
                    className="w-3 h-3 rounded-full border border-gray-300"
                    style={{ backgroundColor: '#CCCCCC' }}
                  />
                )}
                {Array.isArray(safeProduct.colors) && safeProduct.colors.length > 3 && (
                  <span className="text-xs text-gray-500">+{safeProduct.colors.length - 3}</span>
                )}
              </div>
            )}
          </div>
        </div>
      </div>
    </Link>
  );
} 