export class Service {
    static async StartContainer(bucket, grant) {
        alert('starting container: ' + bucket + "/" + grant)
    }
}

export default Service