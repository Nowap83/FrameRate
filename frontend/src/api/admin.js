import apiClient from "./apiClient";

export const adminService = {
    getAllUsers: async () => {
        const response = await apiClient.get('/admin/users');
        return response.data;
    },

    deleteUser: async (id) => {
        const response = await apiClient.delete(`/admin/users/${id}`);
        return response.data;
    },
};
