import {
  Card,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Link } from "react-router-dom";


export default function Login() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-50">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>Welcome to Gazette!</CardTitle>
          <CardDescription>
            Smart RSS aggregator with personalized feeds.
          </CardDescription>
        </CardHeader>
        <CardFooter className="justify-center">
          <Button><Link to="/oauth/login" reloadDocument>Login</Link></Button>
        </CardFooter>
      </Card>
    </div>
  );
}

