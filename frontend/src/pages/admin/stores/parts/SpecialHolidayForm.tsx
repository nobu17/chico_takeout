import * as React from "react";
import { Container, Stack, TextField } from "@mui/material";
import { MobileDatePicker } from "@mui/x-date-pickers/MobileDatePicker";
import { SpecialHoliday } from "../../../../libs/SpecialHoliday";
import { SubmitHandler, useForm, Controller } from "react-hook-form";
import SubmitButtons from "../../../../components/parts/SubmitButtons";

import {
  RequiredErrorMessage,
  MaxLengthErrorMessage,
} from "../../../../libs/ErrorMessages";

import { startDateIsLessThanEndDate } from "../../../../libs/util/TimeCompare";
import { toDate, toDateString } from "../../../../libs/util/DateUtil";

type SpecialHolidayFormProps = {
  editItem: SpecialHoliday;
  onSubmit: callbackSubmit;
  onCancel: callbackCancel;
};
interface callbackSubmit {
  (item: SpecialHoliday): void;
}
interface callbackCancel {
  (): void;
}

type SpecialHolidayInput = {
  id: string;
  name: string;
  start: Date;
  end: Date;
};

const convertInput = (item: SpecialHoliday): SpecialHolidayInput => {
  return {
    id: item.id,
    name: item.name,
    start: toDate(item.start),
    end: toDate(item.end),
  };
};

const reverseInput = (item: SpecialHolidayInput): SpecialHoliday => {
  return {
    id: item.id,
    name: item.name,
    start: toDateString(item.start),
    end: toDateString(item.end),
  };
};

export default function SpecialHolidayForm(props: SpecialHolidayFormProps) {
  const {
    register,
    control,
    handleSubmit,
    setError,
    setValue,
    formState: { errors },
  } = useForm<SpecialHolidayInput>({
    defaultValues: convertInput(props.editItem),
  });
  const onSubmit: SubmitHandler<SpecialHolidayInput> = (data) => {
    const converted = reverseInput(data);
    if (!startDateIsLessThanEndDate(converted.start, converted.end)) {
      setError(`start`, { message: "開始日時が終了時刻よりも大きいです。" });
      setError(`end`, { message: "開始日時が終了時刻よりも大きいです。" });
      return;
    }
    props.onSubmit(converted);
  };
  const onCancel = () => {
    props.onCancel();
  };
  return (
    <>
      <Container maxWidth="sm" sx={{ pt: 5 }}>
        <Stack spacing={3}>
          <TextField
            label="名称"
            {...register("name", {
              required: { value: true, message: RequiredErrorMessage },
              maxLength: { value: 10, message: MaxLengthErrorMessage(10) },
            })}
            error={Boolean(errors.name)}
            helperText={errors.name && errors.name.message}
          />
          <Controller
            name="start"
            control={control}
            render={({ field }) => {
              return (
                <MobileDatePicker
                  {...field}
                  label="開始日"
                  mask="____/__/__ __:__:__"
                  inputFormat="yyyy/MM/dd"
                  renderInput={(props) => (
                    <TextField
                      {...props}
                      error={Boolean(errors.start)}
                      helperText={errors.start && errors.start.message}
                    />
                  )}
                  onChange={(newValue) => {
                    if (newValue != null) {
                      setValue("start", newValue);
                    }
                  }}
                />
              );
            }}
          />
          <Controller
            name="end"
            control={control}
            render={({ field }) => {
              return (
                <MobileDatePicker
                  {...field}
                  label="終了日"
                  mask="____/__/__ __:__:__"
                  inputFormat="yyyy/MM/dd"
                  renderInput={(props) => (
                    <TextField
                      {...props}
                      error={Boolean(errors.end)}
                      helperText={errors.end && errors.end.message}
                    />
                  )}
                  onChange={(newValue) => {
                    if (newValue != null) {
                      setValue("end", newValue);
                    }
                  }}
                />
              );
            }}
          />
          <SubmitButtons
            onSubmit={handleSubmit(onSubmit)}
            onCancel={onCancel}
          />
        </Stack>
      </Container>
    </>
  );
}
