import React from "react";
import { Route, Routes, BrowserRouter } from "react-router-dom";
import {
  createTheme,
  responsiveFontSizes,
  ThemeProvider,
} from "@mui/material/styles";
import { AdapterDateFns } from "@mui/x-date-pickers/AdapterDateFns";
import { LocalizationProvider } from "@mui/x-date-pickers";

import Container from "@mui/material/Container";
import Header from "./components/layouts/Header";
import AdminHome from "./pages/admin/Home";
import ItemKind from "./pages/admin/items/ItemKind";
import StockItem from "./pages/admin/items/StockItem";
import StockItemRemain from "./pages/admin/items/StockItemRemain";
import FoodItem from "./pages/admin/items/FoodItem";
import BusinessHour from "./pages/admin/stores/BusinessHour";
import SpecialBusinessHour from "./pages/admin/stores/SpecialBusinessHour";
import "./App.css";

let theme = createTheme();
theme = responsiveFontSizes(theme);

function App() {
  return (
    <LocalizationProvider dateAdapter={AdapterDateFns}>
      <ThemeProvider theme={theme}>
        <BrowserRouter>
          <Header />
          <Container maxWidth="lg">
            <Routes>
              <Route path="/admin" element={<AdminHome />} />
              <Route path="/admin/items/kind" element={<ItemKind />} />
              <Route path="/admin/items/stock" element={<StockItem />} />
              <Route
                path="/admin/items/stock/remain"
                element={<StockItemRemain />}
              />
              <Route path="/admin/items/food" element={<FoodItem />} />
              <Route path="/admin/store/hour" element={<BusinessHour />} />
              <Route path="/admin/store/special_hour" element={<SpecialBusinessHour />} />
            </Routes>
          </Container>
        </BrowserRouter>
      </ThemeProvider>
    </LocalizationProvider>
  );
}

export default App;
