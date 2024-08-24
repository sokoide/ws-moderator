/** @type {import('next').NextConfig} */
const nextConfig = {
    eslint: {
        // Disable ESLint checks during the build process
        ignoreDuringBuilds: true,
    },
    output: "export", // Configure the app to be exported as static HTML and JavaScript
    trailingSlash: true, // Optional: Adds trailing slashes to paths
};

export default nextConfig;
