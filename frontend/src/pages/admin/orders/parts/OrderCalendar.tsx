import * as React from "react";
import { styled, TextField } from "@mui/material";
import { StaticDatePicker } from "@mui/x-date-pickers";
import { PickersDay, PickersDayProps } from "@mui/x-date-pickers/PickersDay";
import { toDateString } from "../../../../libs/util/DateUtil";
import { UserOrderInfo } from "../../../../libs/apis/order";

type OrderCalendarProps = {
  orders?: UserOrderInfo[];
  onSelected?: callbackOnSelected;
};
interface callbackOnSelected {
  (items: UserOrderInfo[]): void;
}

type CustomPickerDayProps = PickersDayProps<Date> & {
  hasOrder: boolean;
};

const CustomPickersDay = styled(PickersDay, {
  shouldForwardProp: (prop) => prop !== "hasOrder",
})<CustomPickerDayProps>(({ theme, hasOrder }) => ({
  ...(hasOrder && {
    borderRadius: "50%",
    backgroundColor: theme.palette.error.main,
    color: theme.palette.common.white,
    "&:hover, &:focus": {
      backgroundColor: theme.palette.primary.dark,
    },
  }),
})) as React.ComponentType<CustomPickerDayProps>;

export default function OrderCalendar(props: OrderCalendarProps) {
  const [value, setValue] = React.useState<Date | null>(new Date());
  const renderDay = (
    date: Date,
    selectedDates: Array<Date | null>,
    pickersDayProps: PickersDayProps<Date>
  ) => {
    if (!value) {
      return <PickersDay {...pickersDayProps} />;
    }

    const dateStr = toDateString(date);
    let hasOrder = false;
    const foundIndex = props.orders!.findIndex(
      (x) => x.pickupDateTime.startsWith(dateStr)
    );
    if (foundIndex >= 0) {
      hasOrder = true;
    }
    return (
      <CustomPickersDay
        {...pickersDayProps}
        disableMargin
        hasOrder={hasOrder}
      />
    );
  };
  return (
    <>
      <StaticDatePicker
        displayStaticWrapperAs="desktop"
        openTo="day"
        value={value}
        onChange={(newValue) => {
          setValue(newValue);
          if(newValue) {
            const dateStr = toDateString(newValue);
            const sameDates = props.orders!.filter(x => x.pickupDateTime.startsWith(dateStr));
            props.onSelected!(sameDates);
          }
        }}
        renderDay={renderDay}
        renderInput={(params) => <TextField {...params} />}
      />
    </>
  );
}
