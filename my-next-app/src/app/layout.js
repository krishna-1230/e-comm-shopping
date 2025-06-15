import { GeistSans, GeistMono } from 'geist/font'
import "./globals.css";
import Header from "@/components/layout/Header";
import Footer from "@/components/layout/Footer";
import { Toaster } from "react-hot-toast";
import ReduxProvider from "@/redux/provider";

export const metadata = {
  title: "StyleSpace | Modern E-Commerce",
  description: "Your one-stop shop for fashion and accessories",
};

export default function RootLayout({ children }) {
  return (
    <html lang="en" className={`${GeistSans.variable} ${GeistMono.variable}`}>
      <body className="antialiased">
        <ReduxProvider>
          <Toaster position="top-center" />
          <div className="min-h-screen flex flex-col">
            <Header />
            <main className="flex-grow">{children}</main>
            <Footer />
          </div>
        </ReduxProvider>
      </body>
    </html>
  );
}
