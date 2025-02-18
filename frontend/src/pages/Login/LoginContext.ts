import { createContext } from "react";

export type LoginContextType = {
  userName: string;
  userRole: string;
  userToken: string;
  loading: boolean;
  updateUser: (token: string) => void;
};

export const LoginContext = createContext<LoginContextType>({
  userName: "",
  userRole: "",
  userToken: "",
  loading: true,
  updateUser: (_: string) => undefined,
});
