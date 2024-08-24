/** @type {import('next').NextConfig} */
const nextConfig = {
    eslint: {
        // Disable ESLint checks during the build process
        ignoreDuringBuilds: true,
    },
};

export default nextConfig;
