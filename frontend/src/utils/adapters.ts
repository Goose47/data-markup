import { MarkupTypeForm } from "../components/MarkupForm/MarkupForm";
import { MarkupTypeField, MarkupTypeFull } from "./types";

export const markTypeBackendToFrontend = (
  data: MarkupTypeFull
): MarkupTypeForm[] => {
  const mapping: Record<number, MarkupTypeField[]> = {};
  data?.fields.forEach((field) => {
    if (!mapping[field.group_id]) {
      mapping[field.group_id] = [];
    }
    mapping[field.group_id].push(field);
  });
  return Object.values(mapping).map((value) => {
    return {
      type: String(value[0].assessment_type_id),
      label: value[0].label,
      options: value.map((el) => el.name),
    };
  });
};
