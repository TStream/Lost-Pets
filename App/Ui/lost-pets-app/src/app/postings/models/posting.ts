import { IPet, ITag } from "src/app/shared/models/pet";

export interface IPostingsResponse {
    postings: Array<IPosting>
}

export interface IPosting {
    id: number
    date: string
    location: string
    pet: IPet
}

export class PostingRequest {
    constructor(
        public date: string,
        public location: string,
        public pet: IPet
    ) {}

    static adapt(item: any): PostingRequest {
        let tag: ITag = {
            id:item?.tagid,
            shape: item?.tagshape,
            color: item?.tagcolor,
            text: item?.tagtext
        }
        if (item?.pettag){
            tag = item.pettag
        } else if (item?.tag){
            tag = item.tag
        }

        let pet: IPet = {
            id: item?.petid,
            name: item?.petname,
            pictureId: item?.petpictureid,
            color: item?.petcolor,
            marks: item?.petmarks,
            type: item?.pettype,
            typeId: item?.pettypeid,
            breeds: item?.petbreeds,
            tag: tag
        };
        if (item?.pet){
            pet = item.pet
        }


        return new PostingRequest(
            item?.date ? new Date(item.date).toISOString() : "",
            item?.location,
            pet
        );
    }
}