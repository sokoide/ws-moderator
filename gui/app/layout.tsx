import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { AppContextProvider } from "@/context/app-context";
import dynamic from "next/dynamic";

// https://nextjs.org/docs/messages/react-hydration-error
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
            <body className={inter.className} style={{ margin: 0, padding: 0, height: "100vh", display: "flex", flexDirection: "column" }}>
                <AppContextProvider>
                    <Navbar />
                    <main style={{ flex: 1, display: "flex", flexDirection: "column", overflowY: "auto" }}>
                        {children}
                    </main>
                </AppContextProvider>
            </body>
        </html>
    );
}
