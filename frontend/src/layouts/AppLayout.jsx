import Header from "../components/Header";
import Footer from "../components/Footer";

const AppLayout = ({ children }) => {
    return (
        <div className="min-h-screen bg-gray-900 flex flex-col">
            <Header />
            <main className="flex-1 p-8 text-white">
                {children}
            </main>
            <Footer />
        </div>
    );
};

export default AppLayout;