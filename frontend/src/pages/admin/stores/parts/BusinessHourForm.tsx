import * as React from "react";
import { Container, Stack, TextField } from "@mui/material";
import { RhfTimeSelect } from "../../../../components/parts/Rhf/RhfTimeSelect";
import { RhfDayOfWeekSelect } from "../../../../components/parts/Rhf/RhfDayOfWeekSelect";
import { RhfHourOffsetSelect } from "../../../../components/parts/Rhf/RhfHourOffsetSelect";
import { BusinessHour } from "../../../../libs/BusinessHour";
import { StoreTimeList } from "../../../../libs/Constant";
import { SubmitHandler, useForm } from "react-hook-form";
import SubmitButtons from "../../../../components/parts/SubmitButtons";

import { startIsLessThanEnd } from "../../../../libs/util/TimeCompare";
import {
  RequiredErrorMessage,
  MaxLengthErrorMessage,
} from "../../../../libs/ErrorMessages";

type BusinessHourFormProps = {
  editItem: BusinessHour;
  onSubmit: callbackSubmit;
  onCancel: callbackCancel;
};
interface callbackSubmit {
  (itemKind: BusinessHour): void;
}
interface callbackCancel {
  (): void;
}

export default function BusinessHourForm(props: BusinessHourFormProps) {
  const {
    register,
    control,
    handleSubmit,
    setError,
    formState: { errors },
  } = useForm<BusinessHour>({ defaultValues: props.editItem });
  const onSubmit: SubmitHandler<BusinessHour> = (data) => {
    if (!startIsLessThanEnd(data.start, data.end)) {
      setError(`start`, { message: "開始時刻が終了時刻よりも大きいです。" });
      setError(`end`, { message: "開始時刻が終了時刻よりも大きいです。" });
      return;
    }
    props.onSubmit(data);
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
          <RhfTimeSelect
            label="開始時間"
            timeList={timeList}
            name="start"
            control={control}
          />
          <RhfTimeSelect
            label="終了時間"
            timeList={timeList}
            name="end"
            control={control}
          />
          <RhfHourOffsetSelect
            label="予約最終受付時間(X時間前+(1~30分)"
            name="offsetHour"
            control={control}
          />
          <RhfDayOfWeekSelect label="曜日" name="weekdays" control={control} />
          <SubmitButtons
            onSubmit={handleSubmit(onSubmit)}
            onCancel={onCancel}
          />
        </Stack>
      </Container>
    </>
  );
}

const timeList = StoreTimeList;
