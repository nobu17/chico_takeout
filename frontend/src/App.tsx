import React from "react";
import { Route, Routes, BrowserRouter } from "react-router-dom";
import {
  createTheme,
  responsiveFontSizes,
  ThemeProvider,
} from "@mui/material/styles";
import { AdapterDateFns } from "@mui/x-date-pickers/AdapterDateFns";
import ja from "date-fns/locale/ja";
import { LocalizationProvider } from "@mui/x-date-pickers";
import { AdminAuthProvider } from "./components/contexts/AuthContext";

import AdminAuthLayout from "./components/routers/AdminAuthLayout";
import AdminLoginLayout from "./components/routers/AdminLoginLayout";
import UserAuthLayout from "./components/routers/UserAuthLayout";
import UserLoginLayout from "./components/routers/UserLoginLayout";

import Container from "@mui/material/Container";
import Header from "./components/layouts/Header";

import Home from "./pages/Home";
import UserLogin from "./pages/auth/Login";
import UserSignUp from "./pages/auth/SignUp";
import UserReset from "./pages/auth/Reset";
import Logout from "./pages/auth/Logout";

import MyHome from "./pages/mypage/Home";
import ReserveHistory from "./pages/mypage/reserve/ReserveHistory";

import Inquiry from "./pages/Inquiry";

import ReserveHome from "./pages/reserve/Home";

import AdminHome from "./pages/admin/Home";
import AdminLogin from "./pages/admin/Login";

import ItemKind from "./pages/admin/items/ItemKind";
import OptionItem from "./pages/admin/items/OptionItem";
import StockItem from "./pages/admin/items/StockItem";
import StockItemRemain from "./pages/admin/items/StockItemRemain";
import FoodItem from "./pages/admin/items/FoodItem";
import BusinessHour from "./pages/admin/stores/BusinessHour";
import SpecialBusinessHour from "./pages/admin/stores/SpecialBusinessHour";
import SpecialHoliday from "./pages/admin/stores/SpecialHoliday";
import Calendar from "./pages/admin/orders/Calendar";
import AllOrderTable from "./pages/admin/orders/AllOrderTable";
import Monthly from "./pages/admin/statistics/Monthly";
import Messages from "./pages/admin/messages/Messages";

import "./App.css";

let theme = createTheme({
  typography: {
    fontFamily: ["M PLUS Rounded 1c"].join(","),
  },
  palette: {
    primary: {
      light: "#5f5fc4",
      main: "#283593",
      dark: "#001064",
      contrastText: "#ffffff",
    },
  },
});
theme = responsiveFontSizes(theme);

function App() {
  return (
    <LocalizationProvider dateAdapter={AdapterDateFns} adapterLocale={ja}>
      <ThemeProvider theme={theme}>
        <AdminAuthProvider>
          <BrowserRouter>
            <Header />
            <Container maxWidth="lg" disableGutters>
              <Routes>
                <Route element={<UserLoginLayout />}>
                  <Route path="/auth/sign_in" element={<UserLogin />} />
                  <Route path="/auth/sign_up" element={<UserSignUp />} />
                  <Route path="/auth/reset" element={<UserReset />} />
                </Route>
                <Route element={<UserAuthLayout />}>
                  <Route path="/my_page" element={<MyHome />} />
                  <Route path="/my_page/history" element={<ReserveHistory />} />
                  <Route path="/auth/sign_out" element={<Logout />} />
                  <Route path="/reserve" element={<ReserveHome />} />
                </Route>
                <Route element={<AdminLoginLayout />}>
                  <Route path="/admin/sign_in" element={<AdminLogin />} />
                </Route>
                <Route element={<AdminAuthLayout />}>
                  <Route path="/admin" element={<AdminHome />} />
                  <Route path="/admin/items/kind" element={<ItemKind />} />
                  <Route path="/admin/items/option" element={<OptionItem />} />
                  <Route path="/admin/items/stock" element={<StockItem />} />
                  <Route
                    path="/admin/items/stock/remain"
                    element={<StockItemRemain />}
                  />
                  <Route
                    path="/admin/items/stock/remain"
                    element={<StockItemRemain />}
                  />
                  <Route path="/admin/items/food" element={<FoodItem />} />
                  <Route path="/admin/store/hour" element={<BusinessHour />} />
                  <Route
                    path="/admin/store/special_hour"
                    element={<SpecialBusinessHour />}
                  />
                  <Route
                    path="/admin/store/holiday"
                    element={<SpecialHoliday />}
                  />
                  <Route path="/admin/orders/calendar" element={<Calendar />} />
                  <Route
                    path="/admin/orders/all_orders"
                    element={<AllOrderTable />}
                  />
                  <Route
                    path="/admin/statistics/monthly"
                    element={<Monthly />}
                  />
                  <Route
                    path="/admin/messages"
                    element={<Messages />}
                  />
                </Route>
                <Route path="/inquiry" element={<Inquiry />} />
                <Route path="*" element={<Home />} />
              </Routes>
            </Container>
          </BrowserRouter>
        </AdminAuthProvider>
      </ThemeProvider>
    </LocalizationProvider>
  );
}

export default App;
