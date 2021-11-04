import { IPet, ITag } from "src/app/shared/models/pet";

export interface ISightingsResponse {
    sightings: Array<ISighting>
}

export interface ISighting {
    id: number
    date: string
    location: string
    inCustody: boolean
    pet: IPet
}

export class SightingRequest {
    constructor(
        public date: string,
        public location: string,
        inCustody: boolean,
        public pet: IPet
    ) {}

    static adapt(item: any): SightingRequest {
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


        return new SightingRequest(
            item?.date,
            item?.location,
            item?.inCustody,
            pet
        );
    }
}