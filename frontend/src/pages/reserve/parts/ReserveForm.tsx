import * as React from "react";
import {
  Container,
  Stack,
  Stepper,
  StepContent,
  Step,
  StepLabel,
  Box,
  Snackbar,
  Alert,
} from "@mui/material";
import PickupSelect from "./PickupSelect";
import ItemSelect from "./ItemSelect";
import UserInfoInput from "./UserInfoInput";
import ReserveConfirmation from "./ReserveConfirmation";

import { useItemCart } from "../../../hooks/UseItemCart";
import { useUserInfo, UserInfo } from "../../../hooks/UseUserInfo";
import { usePickupDate, PickupDate } from "../../../hooks/UsePickupDate";
import { useOrderableInfo } from "../../../hooks/UseOrderableInfo";
import { useOrder } from "../../../hooks/UseOrder";
import LoadingSpinner from "../../../components/parts/LoadingSpinner";

const steps = ["日時選択", "商品選択", "お客様情報入力", "確認"];

export default function ReserveForm() {
  const [activeStep, setActiveStep] = React.useState(0);
  const [openSnack, setOpenSnack] = React.useState(false);
  const { cart, updateCart, resetCart } = useItemCart();
  const { userInfo, updateUserInfo } = useUserInfo();
  const { pickupDate, updatePickupDate } = usePickupDate();
  const { loading, perDayOrderableInfo, currentOrderableInfo, switchCurrent } =
    useOrderableInfo();
  const { loading: orderLoading, submitOrder } = useOrder();

  const handleNext = () => {
    setActiveStep((prevActiveStep) => prevActiveStep + 1);
  };

  const onPickupDateSubmit = (date: PickupDate) => {
    updatePickupDate(date);
    // switch the item
    switchCurrent(date.date, date.time);
    // reset cart (time is changed equal item is changed.)
    resetCart();

    setActiveStep((prevActiveStep) => prevActiveStep + 1);
  };

  const onUserSubmit = (data: UserInfo) => {
    updateUserInfo(data);
    setActiveStep((prevActiveStep) => prevActiveStep + 1);
  };

  const handleBack = () => {
    setActiveStep((prevActiveStep) => prevActiveStep - 1);
  };

  const handleConfirmSubmit = async () => {
    setOpenSnack(false);
    const result = await submitOrder(pickupDate, cart, userInfo);
    if (result) {
      alert("オーダーしました。");
    } else {
      setOpenSnack(true);
    }
  };

  const displayByStep = (activeStep: number) => {
    if (activeStep === 0) {
      return (
        <PickupSelect
          selectedInfo={pickupDate}
          selectableInfo={perDayOrderableInfo}
          onSubmit={onPickupDateSubmit}
        ></PickupSelect>
      );
    }
    if (activeStep === 1) {
      return (
        <ItemSelect
          allItems={currentOrderableInfo.categories}
          cart={cart}
          onRequestChanged={updateCart}
          onSubmit={handleNext}
          onBack={handleBack}
        ></ItemSelect>
      );
    }
    if (activeStep === 2) {
      return (
        <UserInfoInput
          userInfo={userInfo}
          onSubmit={onUserSubmit}
          onBack={handleBack}
        ></UserInfoInput>
      );
    }
    if (activeStep === 3) {
      return (
        <ReserveConfirmation
          userInfo={userInfo}
          cart={cart}
          onBack={handleBack}
          onSubmit={handleConfirmSubmit}
        ></ReserveConfirmation>
      );
    }
    return <></>;
  };

  return (
    <>
      <Container maxWidth="md" sx={{ pt: 2 }}>
        <Box sx={{ width: "100%" }}>
          <Stepper activeStep={activeStep} orientation="vertical">
            {steps.map((step, index) => (
              <Step key={index}>
                <StepLabel>{step}</StepLabel>
                <StepContent>
                  <Stack spacing={3}>{displayByStep(index)}</Stack>
                </StepContent>
              </Step>
            ))}
          </Stepper>
        </Box>
        <LoadingSpinner
          message="Loading..."
          isLoading={loading || orderLoading}
        />
        <Snackbar
          open={openSnack}
          autoHideDuration={6000}
          anchorOrigin={{ vertical: "top", horizontal: "center" }}
        >
          <Alert severity="error" sx={{ width: "100%" }}>
            オーダー中に問題が発生しました。ご迷惑をおかけしますが、再度お時間を置いてお試しいただくようにお願いいたします。
          </Alert>
        </Snackbar>
      </Container>
    </>
  );
}
