export const getAvatarUrl = (path) => {
    if (!path) return null;
    if (path.startsWith('http')) return path;

    let baseUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

    if (baseUrl.endsWith('/api')) {
        baseUrl = baseUrl.slice(0, -4);
    }

    if (baseUrl.endsWith('/')) {
        baseUrl = baseUrl.slice(0, -1);
    }

    const cleanPath = path.startsWith('/') ? path : `/${path}`;

    return `${baseUrl}${cleanPath}`;
};
