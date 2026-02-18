import Header from "../components/Header";
import Footer from "../components/Footer";

const AppLayout = ({ children }) => {
    return (
        <div className="min-h-screen flex flex-col font-sans text-white">
            <Header />
            <main className="flex-1 flex flex-col relative w-full">
                {children}
            </main>
            <Footer />
        </div>
    );
};

export default AppLayout;