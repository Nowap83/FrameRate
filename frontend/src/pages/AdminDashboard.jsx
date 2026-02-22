import { useState, useEffect } from "react";
import { useAuth } from "../context/AuthContext";
import { Navigate } from "react-router-dom";
import { adminService } from "../api/admin";
import { Users, Trash2, LayoutDashboard, Search, ShieldAlert } from "lucide-react";
import useDocumentTitle from "../hooks/useDocumentTitle";

export default function AdminDashboard() {
    const { user } = useAuth();
    const [users, setUsers] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [activeTab, setActiveTab] = useState("users");
    const [searchQuery, setSearchQuery] = useState("");

    useDocumentTitle("Admin Dashboard");

    // Only admins can access
    if (!user || user.is_admin !== true) {
        return <Navigate to="/" replace />;
    }

    useEffect(() => {
        fetchUsers();
    }, []);

    const fetchUsers = async () => {
        try {
            const data = await adminService.getAllUsers();
            setUsers(data || []);
            setError(null);
        } catch (err) {
            console.error(err);
            setError("Failed to fetch users.");
        } finally {
            setLoading(false);
        }
    };

    const handleDeleteUser = async (id, username) => {
        if (window.confirm(`Are you sure you want to delete user @${username}? This action is irreversible.`)) {
            try {
                await adminService.deleteUser(id);
                setUsers(users.filter(u => u.id !== id));
            } catch (err) {
                console.error("Failed to delete user", err);
                alert("An error occurred while deleting the user.");
            }
        }
    };

    const filteredUsers = users.filter(u =>
        u.username.toLowerCase().includes(searchQuery.toLowerCase()) ||
        u.email.toLowerCase().includes(searchQuery.toLowerCase())
    );

    return (
        <div className="flex flex-col md:flex-row min-h-[calc(100vh-64px)] bg-[var(--color-bg-primary)] text-[var(--color-text-primary)]">
            {/* Lateral Menu (Sidebar) */}
            <aside className="w-full md:w-64 bg-[var(--color-bg-secondary)] border-r border-[#ffffff1a] p-6 flex flex-col gap-4">
                <div className="flex items-center gap-3 text-xl font-bold text-[var(--color-text-primary)] mb-4">
                    <ShieldAlert className="text-[var(--color-accent)]" />
                    Admin Panel
                </div>

                <nav className="flex flex-col gap-2">
                    <button
                        onClick={() => setActiveTab("dashboard")}
                        className={`flex items-center gap-3 px-4 py-3 rounded-lg transition-all ${activeTab === 'dashboard' ? 'bg-[var(--color-accent)] text-white' : 'text-[var(--color-text-secondary)] hover:bg-[#ffffff1a] hover:text-white'}`}
                    >
                        <LayoutDashboard size={20} />
                        Overview
                    </button>
                    <button
                        onClick={() => setActiveTab("users")}
                        className={`flex items-center gap-3 px-4 py-3 rounded-lg transition-all ${activeTab === 'users' ? 'bg-[var(--color-accent)] text-white' : 'text-[var(--color-text-secondary)] hover:bg-[#ffffff1a] hover:text-white'}`}
                    >
                        <Users size={20} />
                        Users
                    </button>
                </nav>
            </aside>

            {/* Main Content Area */}
            <main className="flex-1 p-6 md:p-8">
                <div className="max-w-6xl mx-auto">
                    {activeTab === "dashboard" && (
                        <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
                            <h1 className="text-3xl font-bold heading-font">Admin Dashboard</h1>
                            <p className="text-[var(--color-text-secondary)]">Welcome to the administration panel. Use the sidebar to navigate.</p>

                            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-8">
                                <div className="bg-[var(--color-bg-secondary)] border border-[#ffffff1a] rounded-xl p-6 shadow-lg">
                                    <h3 className="text-lg font-medium text-[var(--color-text-secondary)]">Total Users</h3>
                                    <p className="text-4xl font-bold mt-2">{users.length}</p>
                                </div>
                            </div>
                        </div>
                    )}

                    {activeTab === "users" && (
                        <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
                            <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                                <h1 className="text-3xl font-bold heading-font">Manage Users</h1>
                                <div className="relative">
                                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={18} />
                                    <input
                                        type="text"
                                        placeholder="Search users..."
                                        value={searchQuery}
                                        onChange={(e) => setSearchQuery(e.target.value)}
                                        className="pl-10 pr-4 py-2 bg-[var(--color-bg-secondary)] border border-[#ffffff1a] rounded-lg focus:outline-none focus:border-[var(--color-accent)] transition-colors w-full md:w-64"
                                    />
                                </div>
                            </div>

                            {error && (
                                <div className="p-4 bg-red-500/10 border border-red-500/20 text-red-500 rounded-lg">
                                    {error}
                                </div>
                            )}

                            {loading ? (
                                <div className="flex justify-center py-12">
                                    <div className="loader"></div>
                                </div>
                            ) : (
                                <div className="bg-[var(--color-bg-secondary)] border border-[#ffffff1a] rounded-xl overflow-hidden shadow-lg mt-6">
                                    <div className="overflow-x-auto">
                                        <table className="w-full text-left border-collapse">
                                            <thead>
                                                <tr className="bg-[#ffffff0a] border-b border-[#ffffff1a]">
                                                    <th className="p-4 font-semibold text-sm tracking-wider text-[var(--color-text-secondary)]">User</th>
                                                    <th className="p-4 font-semibold text-sm tracking-wider text-[var(--color-text-secondary)]">Email</th>
                                                    <th className="p-4 font-semibold text-sm tracking-wider text-[var(--color-text-secondary)]">Role</th>
                                                    <th className="p-4 font-semibold text-sm tracking-wider text-[var(--color-text-secondary)] text-right">Actions</th>
                                                </tr>
                                            </thead>
                                            <tbody className="divide-y divide-[#ffffff1a]">
                                                {filteredUsers.length === 0 ? (
                                                    <tr>
                                                        <td colSpan="4" className="p-8 text-center text-[var(--color-text-secondary)]">No users found.</td>
                                                    </tr>
                                                ) : (
                                                    filteredUsers.map((u) => (
                                                        <tr key={u.id} className="hover:bg-[#ffffff0a] transition-colors">
                                                            <td className="p-4">
                                                                <div className="flex items-center gap-3">
                                                                    {u.profile_picture_url ? (
                                                                        <img src={`http://localhost:8080${u.profile_picture_url}`} alt={u.username} className="w-10 h-10 rounded-full object-cover border border-[#ffffff1a]" />
                                                                    ) : (
                                                                        <div className="w-10 h-10 rounded-full bg-[var(--color-accent)]/20 flex items-center justify-center text-[var(--color-accent)] font-bold border border-[var(--color-accent)]/30">
                                                                            {u.username.charAt(0).toUpperCase()}
                                                                        </div>
                                                                    )}
                                                                    <div>
                                                                        <p className="font-medium text-white">@{u.username}</p>
                                                                        <p className="text-xs text-[var(--color-text-secondary)]">ID: {u.id}</p>
                                                                    </div>
                                                                </div>
                                                            </td>
                                                            <td className="p-4 text-[var(--color-text-secondary)]">
                                                                {u.email}
                                                            </td>
                                                            <td className="p-4">
                                                                {u.is_admin ? (
                                                                    <span className="px-2.5 py-1 text-xs font-medium bg-[var(--color-accent)]/20 text-[var(--color-accent)] border border-[var(--color-accent)]/30 rounded-full">Admin</span>
                                                                ) : (
                                                                    <span className="px-2.5 py-1 text-xs font-medium bg-gray-500/20 text-gray-400 border border-gray-500/30 rounded-full">User</span>
                                                                )}
                                                            </td>
                                                            <td className="p-4 text-right">
                                                                {u.id !== user.id && (
                                                                    <button
                                                                        onClick={() => handleDeleteUser(u.id, u.username)}
                                                                        className="p-2 text-red-400 hover:text-red-300 hover:bg-red-400/10 rounded-lg transition-colors"
                                                                        title="Delete user"
                                                                    >
                                                                        <Trash2 size={18} />
                                                                    </button>
                                                                )}
                                                            </td>
                                                        </tr>
                                                    ))
                                                )}
                                            </tbody>
                                        </table>
                                    </div>
                                </div>
                            )}
                        </div>
                    )}
                </div>
            </main>
        </div>
    );
}
