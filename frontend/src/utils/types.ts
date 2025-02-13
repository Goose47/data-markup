export type CreateMarkupTypeFields = {
  name: string;
  assessment_type_id: number;
  group_id: number;
};

export type CreateMarkupTypeRq = {
  name: string;
  fields: CreateMarkupTypeFields[];
};
