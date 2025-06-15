"use client";

import { useState, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import Link from "next/link";
import { FaTrash, FaArrowLeft } from "react-icons/fa";
import toast from "react-hot-toast";
import { fetchCart, updateCartItemQuantity, removeCartItem, clearCartItems } from "@/redux/slices/cartSlice";
import { isAuthenticated } from "@/services/auth-helper";
import { useRouter } from "next/navigation";

export default function CartPage() {
  const dispatch = useDispatch();
  const router = useRouter();
  const { items: cartItems, loading, error, subtotal: cartSubtotal, totalItems } = useSelector(state => state.cart);
  const [promoCode, setPromoCode] = useState("");
  const [promoApplied, setPromoApplied] = useState(false);
  const [promoDiscount, setPromoDiscount] = useState(0);

  useEffect(() => {
    // Check if user is authenticated
    if (!isAuthenticated()) {
      toast.error("Please log in to view your cart");
      router.push('/auth/login');
      return;
    }
    
    // Fetch cart from the API
    dispatch(fetchCart());
  }, [dispatch, router]);

  const handleUpdateQuantity = (id, newQuantity) => {
    if (newQuantity < 1) return;
    
    dispatch(updateCartItemQuantity({
      id,
      updates: { quantity: newQuantity }
    }))
      .unwrap()
      .then(() => {
        toast.success("Cart updated");
      })
      .catch(error => {
        toast.error(error || "Failed to update cart");
      });
  };

  const handleRemoveItem = (id) => {
    dispatch(removeCartItem(id))
      .unwrap()
      .then(() => {
        toast.success("Item removed from cart");
      })
      .catch(error => {
        toast.error(error || "Failed to remove item");
      });
  };

  const handleClearCart = () => {
    dispatch(clearCartItems())
      .unwrap()
      .then(() => {
        toast.success("Cart cleared");
      })
      .catch(error => {
        toast.error(error || "Failed to clear cart");
      });
  };

  const applyPromoCode = (e) => {
    e.preventDefault();
    
    // Simple promo code logic
    if (promoCode.toLowerCase() === "discount10") {
      setPromoApplied(true);
      setPromoDiscount(10);
      toast.success("Promo code applied successfully!");
    } else {
      toast.error("Invalid promo code");
    }
  };

  // Calculate totals
  const subtotal = cartSubtotal || 0;
  const shipping = subtotal > 100 ? 0 : 10;
  const promoAmount = (subtotal * promoDiscount) / 100;
  const total = subtotal + shipping - promoAmount;

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold mb-8">Your Cart</h1>
        <div className="bg-white p-6 rounded-lg shadow-sm border text-center">
          <p className="mb-4 text-xl text-red-500">Error loading your cart</p>
          <p className="mb-8 text-gray-500">{error}</p>
          <button
            onClick={() => dispatch(fetchCart())}
            className="bg-black text-white px-6 py-3 rounded-md hover:bg-gray-800 transition-colors"
          >
            Try Again
          </button>
        </div>
      </div>
    );
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold mb-8">Your Cart</h1>
        <div className="animate-pulse">
          <div className="h-20 bg-gray-200 rounded-md mb-4"></div>
          <div className="h-20 bg-gray-200 rounded-md mb-4"></div>
          <div className="h-40 bg-gray-200 rounded-md"></div>
        </div>
      </div>
    );
  }

  if (!cartItems || cartItems.length === 0) {
    return (
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold mb-8">Your Cart</h1>
        <div className="bg-white p-6 rounded-lg shadow-sm border text-center">
          <p className="mb-4 text-xl">Your cart is empty</p>
          <p className="mb-8 text-gray-500">Looks like you haven't added any products to your cart yet.</p>
          <Link
            href="/products"
            className="bg-black text-white px-6 py-3 rounded-md hover:bg-gray-800 transition-colors"
          >
            Continue Shopping
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Your Cart ({totalItems} {totalItems === 1 ? 'item' : 'items'})</h1>
      
      <div className="flex flex-col lg:flex-row gap-8">
        {/* Cart Items */}
        <div className="lg:w-2/3">
          <div className="bg-white rounded-lg shadow-sm border overflow-hidden">
            <div className="hidden md:grid grid-cols-12 gap-4 p-4 border-b bg-gray-50">
              <div className="col-span-6">
                <span className="font-medium">Product</span>
              </div>
              <div className="col-span-2 text-center">
                <span className="font-medium">Price</span>
              </div>
              <div className="col-span-2 text-center">
                <span className="font-medium">Quantity</span>
              </div>
              <div className="col-span-2 text-right">
                <span className="font-medium">Total</span>
              </div>
            </div>
            
            {/* Cart items list */}
            {cartItems.map(item => {
              const price = item.price || 0;
              const discount = item.discount || 0;
              const itemTotal = (price * (1 - discount / 100)) * item.quantity;
              
              return (
                <div key={item.id} className="grid grid-cols-1 md:grid-cols-12 gap-4 p-4 border-b items-center">
                  {/* Product info */}
                  <div className="col-span-6 flex items-center space-x-4">
                    <div className="w-20 h-20 bg-gray-100 rounded flex-shrink-0">
                      {/* Replace with proper image when available */}
                      <div
                        className="w-full h-full"
                        style={{
                          backgroundImage: `url(${item.image || "/images/placeholder-product.jpg"})`,
                          backgroundSize: "cover",
                          backgroundPosition: "center",
                        }}
                      />
                    </div>
                    <div>
                      <Link
                        href={`/products/${item.slug || item.id}`}
                        className="font-medium hover:text-gray-500 transition-colors"
                      >
                        {item.name}
                      </Link>
                      <p className="text-sm text-gray-500">
                        {item.color && `Color: ${item.color}`}
                        {item.size && item.color && ', '}
                        {item.size && `Size: ${item.size}`}
                      </p>
                    </div>
                  </div>
                  
                  {/* Price */}
                  <div className="md:col-span-2 flex justify-between md:justify-center items-center">
                    <span className="md:hidden font-medium">Price:</span>
                    <span>
                      ${((price * (1 - discount / 100)) || 0).toFixed(2)}
                    </span>
                  </div>
                  
                  {/* Quantity */}
                  <div className="md:col-span-2 flex justify-between md:justify-center items-center">
                    <span className="md:hidden font-medium">Quantity:</span>
                    <div className="flex items-center border rounded overflow-hidden">
                      <button
                        className="px-3 py-1 bg-gray-100 hover:bg-gray-200 transition-colors"
                        onClick={() => handleUpdateQuantity(item.id, item.quantity - 1)}
                        disabled={loading}
                      >
                        -
                      </button>
                      <span className="px-3 py-1">{item.quantity}</span>
                      <button
                        className="px-3 py-1 bg-gray-100 hover:bg-gray-200 transition-colors"
                        onClick={() => handleUpdateQuantity(item.id, item.quantity + 1)}
                        disabled={loading}
                      >
                        +
                      </button>
                    </div>
                  </div>
                  
                  {/* Total */}
                  <div className="md:col-span-2 flex justify-between md:justify-end items-center">
                    <span className="md:hidden font-medium">Total:</span>
                    <div className="flex items-center space-x-2">
                      <span className="font-medium">${itemTotal.toFixed(2)}</span>
                      <button
                        onClick={() => handleRemoveItem(item.id)}
                        className="text-gray-400 hover:text-red-500 transition-colors"
                        disabled={loading}
                      >
                        <FaTrash size={14} />
                      </button>
                    </div>
                  </div>
                </div>
              );
            })}
            
            {/* Clear cart and continue shopping */}
            <div className="flex justify-between p-4">
              <Link
                href="/products"
                className="flex items-center text-gray-600 hover:text-black transition-colors"
              >
                <FaArrowLeft className="mr-2" /> Continue Shopping
              </Link>
              <button
                onClick={handleClearCart}
                className="text-red-500 hover:text-red-700 transition-colors"
                disabled={loading}
              >
                Clear Cart
              </button>
            </div>
          </div>
        </div>
        
        {/* Order Summary */}
        <div className="lg:w-1/3">
          <div className="bg-white rounded-lg shadow-sm border p-6">
            <h2 className="text-xl font-bold mb-4">Order Summary</h2>
            
            <div className="space-y-2 mb-4">
              <div className="flex justify-between">
                <span>Subtotal ({totalItems} items)</span>
                <span>${subtotal.toFixed(2)}</span>
              </div>
              <div className="flex justify-between">
                <span>Shipping</span>
                <span>{shipping === 0 ? "Free" : `$${shipping.toFixed(2)}`}</span>
              </div>
              {promoApplied && (
                <div className="flex justify-between text-green-600">
                  <span>Promo ({promoDiscount}% off)</span>
                  <span>-${promoAmount.toFixed(2)}</span>
                </div>
              )}
            </div>
            
            <div className="border-t pt-4 mb-6">
              <div className="flex justify-between font-bold">
                <span>Total</span>
                <span>${total.toFixed(2)}</span>
              </div>
            </div>
            
            {/* Promo code form */}
            <form onSubmit={applyPromoCode} className="mb-6">
              <p className="text-sm mb-2">Have a promo code?</p>
              <div className="flex space-x-2">
                <input
                  type="text"
                  value={promoCode}
                  onChange={e => setPromoCode(e.target.value)}
                  placeholder="Enter promo code"
                  className="flex-1 p-2 border rounded focus:outline-none focus:ring-2 focus:ring-black"
                  disabled={promoApplied}
                />
                <button
                  type="submit"
                  className={`px-4 py-2 rounded ${
                    promoApplied
                      ? "bg-gray-300 cursor-not-allowed"
                      : "bg-black text-white hover:bg-gray-800"
                  }`}
                  disabled={promoApplied || !promoCode}
                >
                  Apply
                </button>
              </div>
              {promoApplied && (
                <p className="text-sm text-green-600 mt-2">
                  Promo code applied!
                </p>
              )}
            </form>
            
            <Link
              href="/checkout"
              className="block w-full bg-black text-white text-center py-3 rounded-md hover:bg-gray-800 transition-colors"
            >
              Proceed to Checkout
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
} 