import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { AppContextProvider } from "@/context/app-context";
import dynamic from "next/dynamic";

const Navbar = dynamic(() => import("./components/navbar"), { ssr: false });

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
    title: "Family Day",
    description: "AI event",
};

export default function RootLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html lang="en">
            <body className={inter.className}>
                <AppContextProvider>
                    <Navbar />
                    {children}
                </AppContextProvider>
            </body>
        </html>
    );
}
