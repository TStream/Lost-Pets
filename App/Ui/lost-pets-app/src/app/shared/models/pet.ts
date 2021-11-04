export interface IPet {
    id: number
    pictureId: number
    name: string
    color: string
    marks: string
    type: string
    typeId: string
    breeds: Array<string>
    tag: ITag
}

export interface ITag {
    id: number
    shape: string
    color: string
    text: string
}

export interface IPetType {
    id: number
    name: string
}