import Image from "next/image";
import Link from "next/link";
import FeaturedProducts from "@/components/product/FeaturedProducts";
import CategoryGrid from "@/components/product/CategoryGrid";
import HeroSection from "@/components/layout/HeroSection";
import Newsletter from "@/components/layout/Newsletter";

export default function Home() {
  return (
    <div className="min-h-screen">
      {/* Hero Section */}
      <HeroSection />

      {/* Category Grid */}
      <section className="py-12 bg-gray-50">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold text-center mb-8">Shop by Category</h2>
          <CategoryGrid />
        </div>
      </section>

      {/* Featured Products */}
      <section className="py-12">
        <div className="container mx-auto px-4">
          <div className="flex justify-between items-center mb-8">
            <h2 className="text-3xl font-bold">Featured Products</h2>
            <Link 
              href="/products" 
              className="text-gray-600 hover:text-black transition-colors"
            >
              View All →
            </Link>
          </div>
          <FeaturedProducts />
        </div>
      </section>

      {/* Newsletter */}
      <Newsletter />
    </div>
  );
}
