import { useState } from "react";

export type PickupDate = {
  date: string;
  time: string;
};

const defaultPickupDate: PickupDate = {
  date: "",
  time: "",
};

export function usePickupDate() {
  const [pickupDate, setPickupDate] = useState<PickupDate>(defaultPickupDate);
  const updatePickupDate = (request: PickupDate) => {
    setPickupDate({...pickupDate, ...request});
  };

  return {
    pickupDate,
    updatePickupDate,
  };
}
