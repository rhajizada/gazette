import { Footer } from "@/components/Footer";

export default function NotFound() {
  return (
    <div className="flex flex-col h-screen bg-gradient-to-br from-red-50 to-white">
      <div className="flex-1 flex flex-col items-center justify-center p-6 overflow-y-auto">
        <img
          src="/404.gif"
          alt="404"
          className="w-2/3 max-w-sm mb-8 rounded-lg shadow-lg"
        />

        <p className="text-xl text-gray-700 mb-6">
          Oops! This page wandered off into the multiverse.
        </p>
      </div>
      <Footer />
    </div>
  );
}
