import { Card } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Link } from "react-router-dom"
import { Footer } from "@/components/Footer"

export default function Login() {
  return (
    <div className="h-screen flex flex-col justify-between">
      <div className="flex flex-1 items-center justify-center p-4">
        <Card className="w-full max-w-md rounded-xl shadow-xl overflow-hidden">
          <div className="flex flex-col items-center p-6 space-y-6">
            <img src="/logo512.png" alt="Gazette Logo" className="w-20 h-20" />

            <h1 className="text-3xl font-extrabold text-gray-900">
              Welcome to Gazette
            </h1>

            <p className="text-center text-gray-600">
              Smart RSS aggregator with personalized feeds.
            </p>

            <Button asChild size="lg" className="w-full">
              <Link to="/oauth/login" reloadDocument>
                Login
              </Link>
            </Button>
          </div>
        </Card>
      </div>

      <Footer />
    </div>
  )
}

